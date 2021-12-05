package model

import (
	"time"
)

type Member struct {
	GroupID     string    `gorm:"primaryKey" json:"groupID,omitempty"`
	UserID      string    `gorm:"primaryKey" json:"userID,omitempty"`
	MemberSince time.Time `json:"memberSince" faker:"-"`
	IsAdmin     bool      `json:"isAdmin"`
}
