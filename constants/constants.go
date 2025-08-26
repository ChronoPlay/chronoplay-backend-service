package constants

type JsonResp struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
