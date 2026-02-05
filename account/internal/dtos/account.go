package dtos

type ValidateAccountRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
}

type ValidateAccountResponse struct {
	Valid bool `json:"valid"`
}

type DebitRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
	Reason    string `json:"reason" binding:"required"`
}

type DebitResponse struct {
	Balance int64  `json:"balance"`
	Reason  string `json:"reason"`
}

type CreditRequest struct {
	AccountID string `json:"account_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,gt=0"`
	Reason    string `json:"reason" binding:"required"`
}

type CreditResponse struct {
	Balance int64  `json:"balance"`
	Reason  string `json:"reason"`
}
