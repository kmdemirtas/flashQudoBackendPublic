package model

import (
	"time"
)

type Bundle struct {
	ID          string    `gorm:"primaryKey" json:"id" faker:"uuid_digit"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GroupID     string    `json:"groupID"`
	Group       Group     `gorm:"foreignKey:GroupID" json:"-"`
	CreatedAt   time.Time `json:"createdAt" faker:"-"`
	UpdatedAt   time.Time `json:"updatedAt" faker:"-"`
	Cards       []Card    `gorm:"foreignKey:bundle_id" json:"cards,omitempty" faker:"-"`
}
