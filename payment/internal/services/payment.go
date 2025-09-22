package services

import (
	"context"
	"fmt"

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

func (s *PaymentService) CreatePayment(ctx context.Context, param dtos.PaymentRequest) (*dtos.PaymentResponse, error) {
	orderID := uuid.New().String()
	workflowID := fmt.Sprintf("payment-%s", uuid.New().String())
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: workflows.PaymentWorkFlowQueue,
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

	return &dtos.PaymentResponse{WorkflowID: workflowID, RunID: wRun.GetRunID()}, nil
}

func (s *PaymentService) FraudCheck(ctx context.Context, param dtos.FraudCheckRequest) error {
	return nil
}
