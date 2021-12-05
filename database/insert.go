package database

import (
	"time"

	"github.com/ironstone95/FlashQudoV2/generator"
	"github.com/ironstone95/FlashQudoV2/model"
	"gorm.io/gorm"
)

func (db *Database) InsertUser(u model.User) (*model.User, error) {
	if err := db.db.Create(&u).Error; err != nil {
		db.logError("InsertUser", err.Error(), u)
		return nil, ErrGormCreate
	}
	return &u, nil
}

func (db *Database) InsertGroup(g model.Group, userID string) (*model.Group, error) {
	g.ID = generator.CreateID()
	err := db.db.Transaction(func(tx *gorm.DB) error {
		if err := db.db.Create(&g).Error; err != nil {
			db.logError("InsertGroup", err.Error(), g, userID)
			return ErrGormCreate
		}
		m := model.Member{GroupID: g.ID, UserID: userID, IsAdmin: true}
		if _, err := db.InsertMember(m); err != nil {
			return ErrGormCreate
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (db *Database) InsertMember(member model.Member) (*model.Member, error) {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		member.MemberSince = time.Now()
		g := model.Group{ID: member.GroupID}
		u := model.User{ID: member.UserID}
		if err := db.db.Model(&model.User{ID: member.UserID}).Association("Groups").Append(&g); err != nil {
			db.logError("InsertMember", err.Error(), member)
			return ErrGormAssocAppend
		}
		if err := db.db.Model(&model.Group{ID: member.GroupID}).Association("Members").Append(&u); err != nil {
			db.logError("InsertMember", err.Error(), member)
			return ErrGormAssocAppend
		}
		if err := db.db.Save(&member).Error; err != nil {
			db.logError("InsertMember", err.Error(), member)
			return ErrGormSave
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (db *Database) InsertBundle(bundle model.Bundle) (*model.Bundle, error) {
	bundle.ID = generator.CreateID()
	bundle.CreatedAt = time.Now()
	bundle.UpdatedAt = time.Now()
	if err := db.db.Create(&bundle).Error; err != nil {
		db.logError("InsertBundle", err.Error(), bundle)
		return nil, ErrGormCreate
	}
	return &bundle, nil
}

func (db *Database) InsertCard(card model.Card) (*model.Card, error) {
	card.ID = generator.CreateID()
	if err := db.db.Create(&card).Error; err != nil {
		db.logError("InsertCard", err.Error(), card)
		return nil, ErrGormCreate
	}
	return &card, nil
}
