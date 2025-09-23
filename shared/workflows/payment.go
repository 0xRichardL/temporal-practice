package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"

	"github.com/0xRichardL/temporal-practice/shared/activities"
)

var (
	PaymentWorkflowTaskQueue = "payment-tasks"
)

func RegisterPaymentWorkflow(w worker.Worker) {
	w.RegisterWorkflowWithOptions(PaymentWorkFlowDefinition, workflow.RegisterOptions{Name: "PaymentWorkflow"})
}

type PaymentWorkFlowParam struct {
	OrderID   string
	AccountID string
	Amount    int64
}

type PaymentWorkFlowResult struct{}

func PaymentWorkFlowDefinition(ctx workflow.Context, param PaymentWorkFlowParam) (workflowResult *PaymentWorkFlowResult, err error) {
	// Define activity options for account-related activities.
	accountActivityOpts := workflow.ActivityOptions{
		TaskQueue:              activities.AccountActivityTaskQueue,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 3},
		ScheduleToCloseTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
	}
	accountCtx := workflow.WithActivityOptions(ctx, accountActivityOpts)

	// Define activity options for notification-related activities.
	notificationActivityOpts := workflow.ActivityOptions{
		TaskQueue:              activities.NotificationActivityTaskQueue,
		RetryPolicy:            &temporal.RetryPolicy{MaximumAttempts: 3},
		ScheduleToCloseTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
	}
	notificationCtx := workflow.WithActivityOptions(ctx, notificationActivityOpts)

	// WORKFLOW'S STEPS:
	// Step 1: Validate account
	validateAccountParam := activities.ValidateAccountActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	validateAccountActivityResult := activities.ValidateAccountActivityResult{}
	err = workflow.ExecuteActivity(accountCtx, activities.ValidateAccountActivityName, validateAccountParam).Get(ctx, &validateAccountActivityResult)
	if err != nil {
		return nil, err
	}
	if !validateAccountActivityResult.Valid {
		// If validation fails, we return a business-level (Application) error.
		// This error will not be retried by default and clearly indicates a business rule failure.
		return nil, temporal.NewApplicationError("account validation failed", "ValidationFailure")
	}
	// Step 2: Debit account
	debitParam := activities.DebitActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}
	debitActivityResult := activities.DebitActivityResult{}
	err = workflow.ExecuteActivity(accountCtx, activities.DebitActivityName, debitParam).Get(ctx, &debitActivityResult)
	if err != nil {
		return nil, err
	}
	// Step 2.1: Set a SAGA compensation, for executed step.
	// The Defer functions chain will work as SAGA compensation queue.
	defer func() {
		if err != nil {
			// SAGA compensation logic.
			// If the workflow fails after the debit, we must credit the account back.
			// This compensation logic should be very robust. We'll use a new disconnected context
			// with a more aggressive retry policy to ensure the credit succeeds, even if the workflow is cancelled.
			disCtx, _ := workflow.NewDisconnectedContext(ctx)

			compensationCtx := workflow.WithActivityOptions(disCtx, workflow.ActivityOptions{
				TaskQueue:              activities.AccountActivityTaskQueue,
				ScheduleToCloseTimeout: 5 * time.Minute,  // Must be set when RetryPolicy is present.
				StartToCloseTimeout:    20 * time.Second, // Give compensation a bit more time per attempt.
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:    time.Second,
					BackoffCoefficient: 2.0,
					MaximumInterval:    time.Minute,
					MaximumAttempts:    10, // Retry compensation more aggressively.
				},
			})
			comErr := workflow.ExecuteActivity(compensationCtx, activities.CreditActivityName, activities.CreditActivityParam{
				AccountID: param.AccountID,
				Amount:    param.Amount,
			}).Get(compensationCtx, nil)
			err = multierr.Append(err, comErr)
		}
	}()

	// Step 3: Start fraud check child workflow.
	// The main workflow will be pending from here till human fraud check has done.
	var fraudCheckResult FraudCheckWorkflowResult
	fraudCheckCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		WorkflowID: "fraud-check-" + param.OrderID,
		TaskQueue:  FraudCheckWorkflowTaskQueue,
	})
	err = workflow.ExecuteChildWorkflow(fraudCheckCtx, FraudCheckWorkflowDefinition, FraudCheckWorkflowParam{OrderID: param.OrderID}).Get(ctx, &fraudCheckResult)
	if err != nil {
		return nil, err
	}
	if !fraudCheckResult.IsValid {
		// The SAGA compensation will be triggered by this error.
		return nil, temporal.NewApplicationError("fraud check failed", "FraudCheckFailure")
	}

	// Step 4: Send user a notification of the payment.
	err = workflow.ExecuteActivity(notificationCtx, activities.NotifyPaymentActivityName, activities.NotifyPaymentActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentWorkFlowResult{}, nil
}
