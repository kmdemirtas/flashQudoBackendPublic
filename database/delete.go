package database

import (
	"github.com/ironstone95/FlashQudoV2/model"
)

// TODO authorization required
func (db *Database) DeleteGroup(groupID string) error {
	if len(groupID) == 0 {
		return ErrParamNotFound
	}
	if db.memberCount(groupID) > 1 {
		db.logError("DeleteGroup", ErrMembersExists.Error(), groupID)
		return ErrMembersExists
	}
	if err := db.db.Delete(&model.Group{ID: groupID}).Error; err != nil {
		db.logError("DeleteGroup", err.Error(), groupID)
		return ErrGormDelete
	}
	return nil
}

func (db *Database) DeleteMember(groupID, userID string) error {
	if len(groupID) == 0 || len(userID) == 0 {
		return ErrParamNotFound
	}
	m := model.Member{GroupID: groupID, UserID: userID}
	if err := db.db.Delete(&m).Error; err != nil {
		db.logError("DeleteMember", err.Error(), groupID, userID)
		return ErrGormDelete
	}
	return nil
}

func (db *Database) DeleteBundle(bundleID string) error {
	if len(bundleID) == 0 {
		return ErrParamNotFound
	}
	if err := db.db.Delete(&model.Bundle{ID: bundleID}).Error; err != nil {
		db.logError("DeleteBundle", err.Error(), bundleID)
		return ErrGormDelete
	}
	return nil
}

func (db *Database) DeleteCard(cardID string) error {
	if len(cardID) == 0 {
		return ErrParamNotFound
	}
	c := model.Card{ID: cardID}
	if err := db.db.Delete(&c).Error; err != nil {
		db.logError("DeleteCard", err.Error(), cardID)
		return ErrGormDelete
	}
	return nil
}

func (db *Database) DeleteBundleCards(bundleID string) error {
	if len(bundleID) == 0 {
		return ErrParamNotFound
	}

	if err := db.db.Where("bundle_id = ?", bundleID).Delete(&model.Card{}).Error; err != nil {
		return err
	}
	return nil
}
