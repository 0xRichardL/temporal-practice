package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	"github.com/0xRichardL/temporal-practice/payment/internal/dtos"
	"github.com/0xRichardL/temporal-practice/shared/workflows"
)

type PaymentService struct {
	temporalClient client.Client
}

func NewPaymentService(temporalClient client.Client) *PaymentService {
	return &PaymentService{
		temporalClient: temporalClient,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, param dtos.CreatePaymentRequest) (*dtos.CreatePaymentResponse, error) {
	orderID := uuid.New().String()
	workflowID := fmt.Sprintf("payment-%s", uuid.New().String())
	workflowOptions := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                workflows.PaymentWorkflowTaskQueue,
		WorkflowExecutionTimeout: 15 * time.Minute,
		WorkflowRunTimeout:       15 * time.Minute,
		WorkflowTaskTimeout:      10 * time.Second,
	}
	wRun, err := s.temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		workflows.PaymentWorkFlowDefinition,
		workflows.PaymentWorkFlowParam{
			OrderID:   orderID,
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dtos.CreatePaymentResponse{
		OrderID:    orderID,
		WorkflowID: workflowID,
		RunID:      wRun.GetRunID(),
	}, nil
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, dto dtos.GetPaymentStatusRequest) (*dtos.GetPaymentStatusResponse, error) {
	currentStepResult, err := s.temporalClient.QueryWorkflow(ctx, dto.WorkflowID, "", workflows.PaymentWorkflowQueryCurrentStep)
	if err != nil {
		return nil, err
	}
	var step string
	if err := currentStepResult.Get(&step); err != nil {
		return nil, fmt.Errorf("Unable to get payment status: %w", err)
	}
	return &dtos.GetPaymentStatusResponse{
		Status: step,
	}, nil
}
