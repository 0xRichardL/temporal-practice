package workflows

import (
	"errors"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"

	"github.com/0xRichardL/temporal-practice/shared/activities"
)

var (
	PaymentWorkFlowQueue     = "payments"
	PaymentWorkflowTaskQueue = "payment-tasks"
)

func RegisterPaymentWorkFlow(c client.Client) {
	go func() {
		w := worker.New(c, PaymentWorkFlowQueue, worker.Options{})
		w.RegisterWorkflow(PaymentWorkFlowDefinition)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start PaymentWorkflow worker", err)
		}
	}()
}

type PaymentWorkFlowParam struct {
	OrderID   string
	AccountID string
	Amount    int64
}

type PaymentWorkFlowResult struct{}

func PaymentWorkFlowDefinition(ctx workflow.Context, param PaymentWorkFlowParam) (workflowResult *PaymentWorkFlowResult, err error) {
	ao := workflow.ActivityOptions{
		TaskQueue:   PaymentWorkflowTaskQueue,
		RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// WORKFLOW'S STEPS:
	// Step 1: Validate account
	validateAccountParam := activities.ValidateAccountActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	validateAccountActivityResult := activities.ValidateAccountActivityResult{}
	err = workflow.ExecuteActivity(ctx, activities.ValidateAccountActivityName, validateAccountParam).Get(ctx, &validateAccountActivityResult)
	if err != nil {
		return nil, err
	}
	// Step 2: Debit account
	debitParam := activities.DebitActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	debitActivityResult := activities.DebitActivityResult{}
	err = workflow.ExecuteActivity(ctx, activities.DebitActivityName, debitParam).Get(ctx, &debitActivityResult)
	if err != nil {
		return nil, err
	}
	// Step 2.1: Set a SAGA compensation, for executed step.
	// The Defer functions chain will work as SAGA compensation queue.
	defer func() {
		if err != nil {
			comErr := workflow.ExecuteActivity(ctx, activities.CreditActivityName, activities.CreditActivityParam{
				AccountID: param.AccountID,
				Amount:    param.Amount,
			}).Get(ctx, nil)
			multierr.Append(err, comErr)
		}
	}()

	// Step 3: Start fraud check child workflow.
	// The main workflow will be pending from here till human fraud check has done.
	var fraudCheckResult FraudCheckWorkflowResult
	fraudCheckCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: "fraud-check-" + param.OrderID,
		TaskQueue:  FraudCheckTaskQueue,
	})
	err = workflow.ExecuteChildWorkflow(fraudCheckCtx, FraudCheckWorkflowDefinition, FraudCheckWorkflowParam{OrderID: param.OrderID}).Get(ctx, &fraudCheckResult)
	if err != nil {
		return nil, err
	}
	if !fraudCheckResult.IsValid {
		return nil, errors.New("fraud check failed")
	}

	// Step 4: Send user a notification of the payment.
	err = workflow.ExecuteActivity(ctx, activities.NotifyPaymentActivityName, activities.NotifyPaymentActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentWorkFlowResult{}, nil
}
