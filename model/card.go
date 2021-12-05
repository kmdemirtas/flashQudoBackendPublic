package model

import (
	"time"
)

type Card struct {
	ID        string    `gorm:"primaryKey" json:"id" faker:"uuid_digit"`
	BundleID  string    `gorm:"index:ux_card_question,unique;not null" json:"bundleID" faker:"-"`
	Bundle    Bundle    `gorm:"foreignKey:BundleID" json:"-" faker:"-"`
	Question  string    `gorm:"index:ux_card_question,unique;not null" json:"question"`
	Answer    string    `json:"answer"`
	UpdatedAt time.Time `json:"updatedAt" faker:"-"`
	CreatedAt time.Time `json:"createdAt" faker:"-"`
}
