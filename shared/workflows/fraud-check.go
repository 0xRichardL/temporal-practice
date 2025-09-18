package workflows

import (
	"go.temporal.io/sdk/workflow"
)

var (
	FraudCheckTaskQueue  = "fraud-check-tasks"
	FraudCheckSignalName = "fraud-check"
)

type FraudCheckWorkflowParam struct {
	OrderID string
}

type FraudCheckWorkflowResult struct {
	IsValid bool
}

func FraudCheckWorkflowDefinition(ctx workflow.Context, param FraudCheckWorkflowParam) (FraudCheckWorkflowResult, error) {
	var result FraudCheckWorkflowResult
	workflow.GetSignalChannel(ctx, FraudCheckSignalName).Receive(ctx, &result)
	return result, nil
}
