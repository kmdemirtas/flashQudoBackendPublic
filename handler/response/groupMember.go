package response

import "time"

type GroupMember struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	ImageURL    string    `json:"imageUrl"`
	IsAdmin     bool      `json:"isAdmin"`
	MemberSince time.Time `json:"memberSince"`
}
