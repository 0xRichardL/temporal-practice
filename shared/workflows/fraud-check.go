package workflows

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	FraudCheckWorkflowTaskQueue = "fraud-check-tasks"
	FraudCheckSignalName        = "fraud-check"
)

func RegisterFraudCheckWorkflow(w worker.Worker) {
	w.RegisterWorkflowWithOptions(FraudCheckWorkflowDefinition, workflow.RegisterOptions{Name: "FraudCheckWorkflow"})
}

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
