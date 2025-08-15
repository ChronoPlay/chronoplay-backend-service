package dto

type AddCardRequest struct {
	CardNumber      string `json:"card_number"`
	CardName        string `json:"card_name"`
	CardDescription string `json:"card_description"`
	TotalCards      uint32 `json:"total_cards"`
	UserId          uint32 `json:"user_id"`
	UserType        string `json:"user_type"`
}

type GetCardRequest struct {
	CardNumber string `json:"card_number"`
}

type GetCardResponse struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Occupied    uint32 `json:"occupied"`
	Total       uint32 `json:"total"`
	Available   uint32 `json:"available"`
	Creator     uint32 `json:"creator"`
}
