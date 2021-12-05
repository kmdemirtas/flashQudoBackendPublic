package response

type GroupBundle struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	GroupID     string `json:"groupID"`
	CardCount   int    `json:"cardCount"`
}
