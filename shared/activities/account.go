package activities

import (
	"context"

	"go.temporal.io/sdk/worker"
)

var AccountActivityTaskQueue = "account-activity-tasks"

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

type AccountActivities interface {
	Register(w worker.Worker)
	ValidateAccount(ctx context.Context, param ValidateAccountActivityParam) (*ValidateAccountActivityResult, error)
	Debit(ctx context.Context, param DebitActivityParam) (*DebitActivityResult, error)
	Credit(ctx context.Context, param CreditActivityParam) (*CreditActivityResult, error)
}
