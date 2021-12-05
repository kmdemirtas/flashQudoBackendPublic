package model

type Token struct {
	Token   string `gorm:"primaryKey"`
	UserID  string `gorm:"primaryKey;uniqueIndex"`
	Expires int64
}
