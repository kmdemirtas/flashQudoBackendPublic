package response

import "time"

type UserGroup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	MemberSince time.Time `json:"memberSince"`
	IsAdmin     bool      `json:"isAdmin"`
}
