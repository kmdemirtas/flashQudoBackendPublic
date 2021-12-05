package model

import "time"

type User struct {
	ID        string    `gorm:"primaryKey" json:"id" faker:"uuid_digit"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username" faker:"username"`
	CreatedAt time.Time `json:"-" faker:"-"`
	UpdatedAt time.Time `json:"-" faker:"-"`
	ImageURL  string    `json:"imageURL"`
	Groups    []Group   `gorm:"many2many:members" faker:"-" json:"groups,omitempty"`
}
