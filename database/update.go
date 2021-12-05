package database

import (
	"time"

	"github.com/ironstone95/FlashQudoV2/model"
)

// UpdateUser TODO Auth required
func (db *Database) UpdateUser(userID string, updates map[string]interface{}) (*model.User, error) {
	if len(userID) == 0 {
		return nil, ErrParamNotFound
	}
	uv := make(map[string]interface{})
	if imageURL, ok := updates["image_url"]; ok {
		uv["image_url"] = imageURL
	} else {
		db.logError("UpdateUser", ErrUpdateValueNotFound.Error(), userID, updates)
		return nil, ErrUpdateValueNotFound
	}

	if err := db.db.Model(&model.User{}).Where("id = ?", userID).Updates(uv).Error; err != nil {
		db.logError("UpdateUser", err.Error(), userID, updates)
		return nil, ErrGormUpdate
	}

	u := model.User{}
	if err := db.db.Where("id = ?", userID).First(&u).Error; err != nil {
		db.logError("UpdateUser", err.Error(), userID, updates)
		return nil, ErrGormGet
	}
	return &u, nil
}

func (db *Database) UpdateBundle(bundleID string, updates map[string]interface{}) (*model.Bundle, error) {
	if len(bundleID) == 0 {
		return nil, ErrParamNotFound
	}
	uv := make(map[string]interface{})
	canUpdate := false
	if title, ok := updates["title"]; ok {
		uv["title"] = title
		canUpdate = true
	}
	if desc, ok := updates["description"]; ok {
		uv["description"] = desc
		canUpdate = true
	}

	if !canUpdate {
		db.logError("UpdateBundle", ErrUpdateValueNotFound.Error(), bundleID, updates)
		return nil, ErrUpdateValueNotFound
	}
	uv["updated_at"] = time.Now()
	b := model.Bundle{ID: bundleID}
	if err := db.db.Model(&b).Updates(uv).Error; err != nil {
		db.logError("UpdateBundle", err.Error(), bundleID, updates)
		return nil, ErrGormUpdate
	}
	bDB := model.Bundle{ID: bundleID}
	if err := db.db.Model(&bDB).First(&bDB).Error; err != nil {
		db.logError("UpdateBundle", err.Error(), bundleID, updates)
		return nil, ErrGormGet
	}
	return &bDB, nil
}

func (db *Database) UpdateCard(cardID string, updates map[string]interface{}) (*model.Card, error) {
	if len(cardID) == 0 {
		return nil, ErrParamNotFound
	}
	canUpdate := false
	uv := make(map[string]interface{})
	if question, ok := updates["question"]; ok {
		uv["question"] = question
		canUpdate = true
	}

	if answer, ok := updates["answer"]; ok {
		uv["answer"] = answer
		canUpdate = true
	}

	if !canUpdate {
		db.logError("UpdateCard", ErrUpdateValueNotFound.Error(), cardID, updates)
		return nil, ErrUpdateValueNotFound
	}

	uv["updated_at"] = time.Now()

	c := model.Card{ID: cardID}
	if err := db.db.Model(&c).Updates(updates).Error; err != nil {
		db.logError("UpdateCard", err.Error(), cardID, updates)
		return nil, ErrGormUpdate
	}
	if err := db.db.Where("id = ?", cardID).First(&c).Error; err != nil {
		db.logError("UpdateCard", err.Error(), cardID, updates)
		return nil, ErrGormGet
	}
	return &c, nil
}

func (db *Database) UpdateMember(groupID, userID string, updates map[string]interface{}) (*model.Member, error) {
	if len(groupID) == 0 || len(userID) == 0 {
		return nil, ErrParamNotFound
	}
	uv := make(map[string]interface{})
	if isAdmin, ok := updates["is_admin"]; !ok {
		db.logError("UpdateMember", ErrUpdateValueNotFound.Error(), groupID, userID, updates)
		return nil, ErrUpdateValueNotFound
	} else {
		uv["is_admin"] = isAdmin
	}
	m := model.Member{GroupID: groupID, UserID: userID}
	if err := db.db.Model(&m).Updates(uv).Error; err != nil {
		db.logError("UpdateMember", err.Error(), groupID, userID, updates)
		return nil, err
	}
	if err := db.db.Where("group_id = ? and user_id = ?", groupID, userID).First(&m).Error; err != nil {
		db.logError("UpdateMember", err.Error(), groupID, userID, updates)
		return nil, ErrGormGet
	}
	return &m, nil
}

func (db *Database) UpdateGroup(groupID string, updates map[string]interface{}) (*model.Group, error) {
	if len(groupID) == 0 {
		return nil, ErrParamNotFound
	}
	uv := make(map[string]interface{})
	if name, ok := updates["name"]; !ok {
		db.logError("UpdateGroup", ErrUpdateValueNotFound.Error(), groupID, updates)
		return nil, ErrUpdateValueNotFound
	} else {
		uv["name"] = name
	}

	g := model.Group{ID: groupID}
	if err := db.db.Model(&g).Updates(uv).Error; err != nil {
		db.logError("UpdateGroup", err.Error(), groupID, updates)
		return nil, ErrGormUpdate
	}

	if err := db.db.Where("ID = ?", groupID).First(&g).Error; err != nil {
		db.logError("UpdateGroup", err.Error(), groupID, updates)
		return nil, ErrGormGet
	}
	return &g, nil
}
