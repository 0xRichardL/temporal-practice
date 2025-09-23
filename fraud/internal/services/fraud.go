package services

import (
	"context"

	"github.com/0xRichardL/temporal-practice/fraud/internal/dtos"
	"github.com/0xRichardL/temporal-practice/shared/workflows"
	"go.temporal.io/sdk/client"
)

type FraudService struct {
	temporalClient client.Client
}

func NewFraudService(temporalClient client.Client) *FraudService {
	return &FraudService{
		temporalClient: temporalClient,
	}
}

func (s *FraudService) Check(ctx context.Context, dto dtos.FraudCheckRequest) error {
	return s.temporalClient.SignalWorkflow(ctx, "fraud-check-"+dto.OrderID, "", workflows.FraudCheckSignalName, workflows.FraudCheckWorkflowResult{
		IsValid: dto.IsValid,
	})
}
