package dto

type AddCardRequest struct {
	CardNumber      string `json:"card_number"`
	CardName        string `json:"card_name"`
	CardDescription string `json:"card_description"`
	TotalCards      uint32 `json:"total_cards"`
	UserId          uint32 `json:"user_id"`
	UserType        string `json:"user_type"`
}
