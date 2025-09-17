package activities

var ValidateAccountActivityName = "account::validate"

type ValidateAccountActivityParam struct {
	AccountID string
	Amount    int64
}

type ValidateAccountActivityResult struct {
	Valid bool
}

var DebitActivityName = "account::debit"

type DebitActivityParam struct {
	AccountID string
	Amount    int64
}

type DebitActivityResult struct {
	Balance int64
}

var CreditActivityName = "account::credit"

type CreditActivityParam struct {
	AccountID string
	Amount    int64
}

type CreditActivityResult struct {
	Balance int64
}
