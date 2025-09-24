package dtos

type CreatePaymentRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type CreatePaymentResponse struct {
	OrderID    string `json:"order_id"`
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
}

type GetPaymentStatusRequest struct {
	WorkflowID string `uri:"workflowID" binding:"required"`
}

type GetPaymentStatusResponse struct {
	Status string `json:"status"`
}
