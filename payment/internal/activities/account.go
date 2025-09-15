package activities

var ValidateAccountActivityName = "account::validate"

type ValidateAccountActivityParam struct {
	AccountID string
	Amount    int64
}

type ValidateAccountActivityResultObject struct {
	Valid bool
}

var DebitActivityName = "account::debit"

type DebitActivityParam struct {
	AccountID string
	Amount    int64
}

type DebitActivityResultObject struct {
	AccountID string
	Balance   int64
}
