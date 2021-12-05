package request

import "github.com/ironstone95/FlashQudoV2/model"

type BundlePostRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"` // CAN BE NULL
}

type CardPostRequest struct {
	Question *string `json:"question"`
	Answer   *string `json:"answer"`
}
type GroupPostRequest struct {
	Name *string `json:"name"`
}
type MemberPostRequest struct {
	GroupID  *string `json:"groupID"`
	Username *string `json:"username"`
	IsAdmin  *bool   `json:"isAdmin"`
}
type UserPostRequest struct {
	Username *string `json:"username"`
	ImageURL *string `json:"imageURL"`
}

// CreateBundle creates bundle if all fields are valid.
func (bpr *BundlePostRequest) CreateBundle() (model.Bundle, error) {
	if bpr.Title == nil {
		return model.Bundle{}, ErrMissingField
	}
	b := model.Bundle{}
	b.Title = *bpr.Title
	if bpr.Description != nil {
		b.Description = *bpr.Description
	}
	return b, nil
}

func (cpr *CardPostRequest) CreateCard() (model.Card, error) {
	if cpr.Question == nil || cpr.Answer == nil {
		return model.Card{}, ErrMissingField
	}

	return model.Card{Question: *cpr.Question, Answer: *cpr.Answer}, nil
}

func (gpr *GroupPostRequest) CreateGroup() (model.Group, error) {
	if gpr.Name == nil {
		return model.Group{}, ErrMissingField
	}
	return model.Group{Name: *gpr.Name}, nil
}

func (mpr *MemberPostRequest) CreateMember(userID string) (model.Member, error) {
	if mpr.GroupID == nil || mpr.Username == nil || mpr.IsAdmin == nil {
		return model.Member{}, ErrMissingField
	}

	return model.Member{
		GroupID: *mpr.GroupID,
		UserID:  userID,
		IsAdmin: *mpr.IsAdmin,
	}, nil
}

func (urp *UserPostRequest) CreateUser() (model.User, error) {
	if urp.Username == nil || urp.ImageURL == nil {
		return model.User{}, ErrMissingField
	}
	return model.User{Username: *urp.Username, ImageURL: *urp.ImageURL}, nil
}
