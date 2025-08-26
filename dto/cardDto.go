package dto

import "mime/multipart"

type AddCardRequest struct {
	CardNumber      string                `json:"card_number" form:"card_number"`
	CardName        string                `json:"card_name" form:"card_name"`
	CardDescription string                `json:"card_description" form:"card_description"`
	TotalCards      uint32                `json:"total_cards" form:"total_cards"`
	UserId          uint32                `json:"user_id" form:"user_id"`
	UserType        string                `json:"user_type" form:"user_type"`
	Image           *multipart.FileHeader `json:"-" form:"image"` // Optional field for image upload
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
	ImageUrl    string `json:"image_url"`
}
