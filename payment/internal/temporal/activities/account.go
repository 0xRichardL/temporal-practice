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
	Balance int64
}

var CreditActivityName = "account::credit"

type CreditActivityParam struct {
	AccountID string
	Amount    int64
}

type CreditActivityResultObject struct {
	Balance int64
}
