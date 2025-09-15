package workflows

import (
	"github.com/0xRichardL/temporal-practice/payment/internal/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type PaymentWorkFlowParam struct {
	AccountID string
	Amount    int64
}

type PaymentWorkFlowResult struct{}

func PaymentWorkFlowDefinition(ctx workflow.Context, param PaymentWorkFlowParam) (*PaymentWorkFlowResult, error) {
	ao := workflow.ActivityOptions{
		TaskQueue:   "payment-tasks",
		RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// 1. Validate account
	validateAccountParam := activities.ValidateAccountActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	validateAccountActivityResultObject := activities.ValidateAccountActivityResultObject{}
	err := workflow.ExecuteActivity(ctx, "account::validate", validateAccountParam).Get(ctx, &validateAccountActivityResultObject)
	if err != nil {
		return nil, err
	}
	// 2. Debit account
	debitParam := activities.DebitActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	debitActivityResultObject := activities.DebitActivityResultObject{}
	err = workflow.ExecuteActivity(ctx, "account::debit", debitParam).Get(ctx, &debitActivityResultObject)
	if err != nil {
		return nil, err
	}

	return &PaymentWorkFlowResult{}, nil
}
