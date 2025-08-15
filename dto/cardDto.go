package dto

type AddCardRequest struct {
	CardNumber      string `json:"card_number"`
	CardDescription string `json:"card_description"`
	TotalCards      uint32 `json:"total_cards"`
}
