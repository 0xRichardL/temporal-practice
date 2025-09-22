package workflows

import (
	"log"
	"time"

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

// RegisterPaymentWorkFlow sets up and starts a worker for the payment workflow.
// In a larger application, it's a good practice to centralize all worker registration
// and startup logic in a dedicated package or in the main application entry point (e.g., cmd/worker/main.go),
// rather than coupling it with the workflow definition file.
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
	// It's a best practice to set a short StartToCloseTimeout that reflects the expected execution time
	// of a single attempt, and a longer ScheduleToCloseTimeout that accounts for queue time and retries.
	// Setting them to the same value can cause timeouts if the activity waits in the task queue.
	ao := workflow.ActivityOptions{
		TaskQueue:   PaymentWorkflowTaskQueue,
		RetryPolicy: &temporal.RetryPolicy{MaximumAttempts: 3},
		// Total time from scheduling to completion, including retries and queue time.
		ScheduleToCloseTimeout: 1 * time.Minute,
		// Max time for a single activity execution attempt.
		StartToCloseTimeout: 10 * time.Second,
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
	err = workflow.ExecuteActivity(ctx, activities.DebitActivityName, debitParam).Get(ctx, &debitActivityResult)
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
				TaskQueue:           PaymentWorkflowTaskQueue,
				StartToCloseTimeout: 20 * time.Second, // Give compensation a bit more time per attempt.
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
		TaskQueue:  FraudCheckTaskQueue,
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
	err = workflow.ExecuteActivity(ctx, activities.NotifyPaymentActivityName, activities.NotifyPaymentActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentWorkFlowResult{}, nil
}
