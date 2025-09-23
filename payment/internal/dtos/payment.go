package dtos

type PaymentRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type PaymentResponse struct {
	OrderID    string `json:"order_id"`
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
}

type FraudCheckRequest struct {
	AccountID string `json:"account_id"`
	PaymentID string `json:"payment_id"`
	IsValid   bool   `json:"is_valid"`
}
