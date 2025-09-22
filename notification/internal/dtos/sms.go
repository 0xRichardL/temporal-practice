package dtos

type SendSMSRequest struct {
	AccountID string `json:"account_id"`
	Message   string `json:"message"`
}
