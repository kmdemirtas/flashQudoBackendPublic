package request

type BundlePatchRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CardPatchRequest struct {
	Question *string `json:"question"`
	Answer   *string `json:"answer"`
}

type UserPatchRequest struct {
	ImageURL *string `json:"imageURL"`
}

type GroupPatchRequest struct {
	Name *string `json:"name"`
}

type MemberPatchRequest struct {
	IsAdmin *bool `json:"isAdmin"`
}

func (bpr *BundlePatchRequest) GetPatchValues() (map[string]interface{}, error) {
	if bpr.Title == nil && bpr.Description == nil {
		return nil, ErrMissingField
	}

	pv := make(map[string]interface{})
	if bpr.Title != nil {
		pv["title"] = *bpr.Title
	}
	if bpr.Description != nil {
		pv["description"] = *bpr.Description
	}
	return pv, nil
}

func (cpr *CardPatchRequest) GetPatchValues() (map[string]interface{}, error) {
	if cpr.Question == nil && cpr.Answer == nil {
		return nil, ErrMissingField
	}

	pv := make(map[string]interface{})
	if cpr.Question == nil {
		pv["answer"] = *cpr.Answer
	} else if cpr.Answer == nil {
		pv["question"] = *cpr.Question
	} else {
		pv["answer"] = *cpr.Answer
		pv["question"] = *cpr.Question
	}
	return pv, nil
}

func (gpr *GroupPatchRequest) GetPatchValues() (map[string]interface{}, error) {
	if gpr.Name == nil {
		return nil, ErrMissingField
	}
	pv := make(map[string]interface{})
	pv["name"] = *gpr.Name
	return pv, nil
}

func (upr *UserPatchRequest) GetPatchValues() (map[string]interface{}, error) {
	if upr.ImageURL == nil {
		return nil, ErrMissingField
	}
	pv := make(map[string]interface{})
	pv["image_url"] = *upr.ImageURL
	return pv, nil
}

func (mpr *MemberPatchRequest) GetPatchValues() (map[string]interface{}, error) {
	if mpr.IsAdmin == nil {
		return nil, ErrMissingField
	}

	pv := make(map[string]interface{})
	pv["is_admin"] = *mpr.IsAdmin
	return pv, nil
}
