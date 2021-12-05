package database

import (
	"time"

	"github.com/ironstone95/FlashQudoV2/model"
)

// AuthToken checks if the token is in the database and token expiration date is not reached.
// If token is expired, then deletes it from the database.
func (db *Database) AuthToken(token string) error {
	t := &model.Token{Token: token}
	if err := db.db.Model(&t).First(&t).Error; err != nil {
		db.logError("AuthToken", err.Error(), token)
		return ErrTokenNotFound
	}
	if t.Expires < time.Now().Unix() {
		if err := db.db.Delete(&t).Error; err != nil {
			db.logError("AuthToken", err.Error(), token)
			return ErrGormDelete
		}
		db.logError("AuthToken", ErrTokenExpired.Error(), token)
		return ErrTokenExpired
	}
	return nil
}

// InserToken inserts given token to the database and returns created token.
func (db *Database) InsertToken(token model.Token) (*model.Token, error) {
	if err := db.db.Create(&token).Error; err != nil {
		db.logError("InsertToken", err.Error(), token.Token, token.UserID)
		return nil, ErrGormCreate
	}
	if err := db.db.Model(&token).First(&token).Error; err != nil {
		db.logError("InsertToken", err.Error(), token.Token, token.UserID)
		return nil, ErrGormGet
	}
	return &token, nil
}

// GetIDFromUsername returns the id of the username.
func (db *Database) GetIDFromUsername(username string) (string, error) {
	if len(username) == 0 {
		return "", ErrParamNotFound
	}
	var userID string
	if err := db.db.Model(&model.User{}).Select("id").Where("username = ?", username).First(&userID).Error; err != nil {
		return "", ErrUserNotFound
	}
	return userID, nil
}

// GetIDFromToken returns the id of the token.
func (db *Database) GetIDFromToken(token string) (string, error) {
	if len(token) == 0 {
		db.logError("GetIDFromToken", ErrParamNotFound.Error(), token)
		return "", ErrParamNotFound
	}
	var id string
	if err := db.db.Model(&model.Token{}).Where("token = ?", token).Select("user_id").First(&id).Error; err != nil {
		db.logError("GetIDFromToken", err.Error(), token)
		return "", ErrGormGet
	}
	if len(id) == 0 {
		db.logError("GetIDFromToken", "", token)
		return "", ErrIDLength
	}
	return id, nil
}

func (db *Database) IsGroupMember(groupID, userID string) (bool, error) {
	if len(groupID) == 0 || len(userID) == 0 {
		db.logError("IsGroupMember", ErrParamNotFound.Error(), groupID, userID)
		return false, ErrParamNotFound
	}
	m := model.Member{UserID: userID, GroupID: groupID}
	if err := db.db.Model(&m).First(&m).Error; err != nil {
		db.logError("IsGroupMember", err.Error(), groupID, userID)
		return false, err
	}
	return true, nil
}

func (db *Database) CanSeeBundle(bundleID, userID string) (bool, error) {
	if len(bundleID) == 0 || len(userID) == 0 {
		db.logError("CanSeeBundle", ErrParamNotFound.Error(), bundleID, userID)
		return false, ErrParamNotFound
	}
	var count int64
	if err := db.db.Model(&model.Member{}).
		Distinct("user_id").
		Joins("left join bundles on members.group_id = bundles.group_id").
		Where("bundles.id = ? and members.user_id = ?", bundleID, userID).
		Count(&count).Error; err != nil {
		db.logError("CanSeeBundle", err.Error(), bundleID, userID)
		return false, ErrGormGet
	}
	return count >= 1, nil
}

func (db *Database) CanEditBundle(bundleID, userID string) (bool, error) {
	if len(bundleID) == 0 || len(userID) == 0 {
		return false, ErrParamNotFound
	}
	var editBundle bool
	if err := db.db.Table("members").Select("members.is_admin").
		Joins("left join bundles on bundles.group_id = members.group_id").
		Where("bundles.id = ? and members.user_id = ?", bundleID, userID).
		Find(&editBundle).Error; err != nil {
		db.logError("CanEditBundle", err.Error(), bundleID, userID)
		return false, ErrGormGet
	}
	return editBundle, nil
}

func (db *Database) CanEditCard(cardID, userID string) (bool, error) {
	if len(cardID) == 0 || len(userID) == 0 {
		return false, ErrParamNotFound
	}
	var isAdmin bool
	if err := db.db.Table("members").Select("members.is_admin").
		Joins("left join bundles on bundles.group_id = members.group_id").
		Joins("left join cards on cards.bundle_id = bundles.id").
		Where("cards.id = ? and members.user_id = ?", cardID, userID).Find(&isAdmin).Error; err != nil {
		db.logError("CanEditCard", err.Error(), cardID, userID)
		return false, ErrGormGet
	}
	return isAdmin, nil
}

func (db *Database) IsAdmin(groupID, userID string) (bool, error) {
	if len(groupID) == 0 || len(userID) == 0 {
		return false, ErrParamNotFound
	}
	var admin bool
	var memberCount int64
	if err := db.db.Table("members").Select("members.is_admin").Where("group_id = ? and user_id = ?", groupID, userID).Find(&admin).Error; err != nil {
		db.logError("IsAdmin", err.Error(), groupID, userID)
		return false, ErrGormGet
	}

	if !admin {
		if err := db.db.Table("members").Where("group_id = ?", groupID).Count(&memberCount).Error; err != nil {
			db.logError("IsAdmin", err.Error(), groupID, userID)
			return false, ErrGormGet
		}
		if memberCount == 0 { // shouldn't enter here
			return true, nil
		}
	}
	return admin, nil
}

func (db *Database) memberCount(groupID string) int64 {
	var count int64
	if err := db.db.Model(&model.Member{GroupID: groupID}).Count(&count).Error; err != nil {
		db.logError("memberCount", err.Error(), groupID)
		return -1
	}
	return count
}
