package workflows_test

import (
	"github.com/0xRichardL/temporal-practice/shared/workflows"
)

func (s *UnitTestSuite) TestFraudCheckWorkflowDefinition_validSignalReceived() {
	param := workflows.FraudCheckWorkflowParam{
		OrderID: "order-1",
	}
	signalValue := workflows.FraudCheckWorkflowResult{
		IsValid: true,
	}
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(workflows.FraudCheckSignalName, signalValue)
	}, 0)
	s.env.ExecuteWorkflow(workflows.FraudCheckWorkflowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())

	var result workflows.FraudCheckWorkflowResult
	s.NoError(s.env.GetWorkflowResult(&result))
	s.True(result.IsValid)
}

func (s *UnitTestSuite) TestFraudCheckWorkflowDefinition_invalidSignalReceived() {
	param := workflows.FraudCheckWorkflowParam{
		OrderID: "order-2",
	}
	signalValue := workflows.FraudCheckWorkflowResult{
		IsValid: false,
	}
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(workflows.FraudCheckSignalName, signalValue)
	}, 0)
	s.env.ExecuteWorkflow(workflows.FraudCheckWorkflowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())

	var result workflows.FraudCheckWorkflowResult
	s.NoError(s.env.GetWorkflowResult(&result))
	s.False(result.IsValid)
}
