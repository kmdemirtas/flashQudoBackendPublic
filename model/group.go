package model

import (
	"time"
)

type Group struct {
	ID        string    `gorm:"primaryKey" json:"id" faker:"uuid_digit"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt" faker:"-"`
	Members   []User    `gorm:"many2many:members" json:"members,omitempty" faker:"-"`
	Bundles   []Bundle  `gorm:"foreignKey:group_id" json:"bundles,omitempty" faker:"-"`
}
