package dtos

type FraudCheckRequest struct {
	OrderID string `json:"order_id,required"`
	IsValid bool   `json:"is_valid"`
}
