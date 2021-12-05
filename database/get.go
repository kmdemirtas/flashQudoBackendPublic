package database

import (
	"github.com/ironstone95/FlashQudoV2/handler/response"
	"github.com/ironstone95/FlashQudoV2/model"
)

// GetUser fetches the first user with given query parameters. Suggested values are id and username.
func (db *Database) GetUser(query map[string]interface{}) (*model.User, error) {
	user := &model.User{}
	if err := db.db.Where(query).First(&user).Error; err != nil {
		db.logError("GetUser", err.Error(), query)
		return nil, ErrGormGet
	}
	return user, nil
}

// GetGroup fetch the first group with given query parameters. Suggested values are ID and Name.
func (db *Database) GetGroup(query map[string]interface{}) (*model.Group, error) {
	group := &model.Group{}
	if err := db.db.Where(query).First(&group).Error; err != nil {
		db.logError("GetGroup", err.Error(), query)
		return nil, ErrGormGet
	}
	return group, nil
}

func (db *Database) GetGroupUsers(groupID string, page, limit int) ([]model.User, error) {
	var users []model.User
	if err := db.db.Model(&model.Group{ID: groupID}).
		Limit(limit).
		Offset(limit * (page - 1)).
		Association("Members").
		Find(&users); err != nil {
		db.logError("GetGroupUsers", err.Error(), groupID, page, limit)
		return nil, err
	}
	return users, nil
}

func (db *Database) GetGroupMembers(groupID string, page, limit int) ([]response.GroupMember, error) {
	var members []response.GroupMember
	if err := db.db.Table("users").
		Select("users.id, users.username, users.image_url, members.member_since, members.is_admin").
		Joins("left join members on users.id = members.user_id").
		Where("members.group_id = ?", groupID).Limit(limit).Offset((page - 1) * limit).Scan(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (db *Database) GetUserGroups(userID string, page, limit int) ([]response.UserGroup, error) {
	var groups []response.UserGroup
	if err := db.db.Table("groups").
		Select("groups.id, groups.name, members.member_since, members.is_admin").
		Joins("left join members on groups.id = members.group_id").
		Where("members.user_id = ?", userID).
		Limit(limit).Offset((page - 1) * limit).Scan(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (db *Database) GetGroupBundles(groupID string, page, limit int) ([]response.GroupBundle, error) {
	var bundles []response.GroupBundle
	if err := db.db.Table("bundles").
		Select("bundles.id, bundles.title, bundles.description, bundles.group_id, count(cards.id) as \"card_count\"").
		Joins("left join cards on cards.bundle_id = bundles.id").Where("bundles.group_id = ?", groupID).
		Group("bundles.id").Limit(limit).Offset((page - 1) * limit).Scan(&bundles).Error; err != nil {
		return nil, err
	}
	return bundles, nil
}

func (db *Database) GetBundle(bundleID string) (*response.GroupBundle, error) {
	b := response.GroupBundle{}
	if err := db.db.Table("bundles").
		Select("bundles.id, bundles.title, bundles.description, bundles.group_id, count(cards.id) as \"card_count\"").
		Joins("left join cards on cards.bundle_id = bundles.id").Where("bundles.id = ?", bundleID).
		Group("bundles.id").Scan(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (db *Database) GetBundleCards(bundleID string, page, limit int) ([]model.Card, error) {
	var cards []model.Card
	if err := db.db.Model(&model.Bundle{ID: bundleID}).
		Limit(limit).
		Offset(limit * (page - 1)).
		Association("Cards").
		Find(&cards); err != nil {
		db.logError("GetBundleCards", err.Error(), bundleID, page, limit)
		return nil, ErrGormGet
	}
	return cards, nil
}

// func (db *Database) getUserID(username string) (string, error) {
// 	var id string
// 	if err := db.db.Model(&model.User{Username: username}).Select("id").First(&id).Error; err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }
