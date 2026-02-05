package workflows_test

import (
	"errors"
	"time"

	"github.com/0xRichardL/temporal-practice/shared/activities"
	"github.com/0xRichardL/temporal-practice/shared/workflows"
	"github.com/stretchr/testify/mock"
)

func (s *UnitTestSuite) TestPaymentWorkFlowDefinition_InvalidAccount() {
	param := workflows.PaymentWorkFlowParam{
		OrderID:   "order-1",
		AccountID: "account-1",
		Amount:    100,
	}
	validateAccountParam := activities.ValidateAccountActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	}

	// Mock the activity with expected behavior
	s.env.OnActivity(activities.ValidateAccountActivityName, validateAccountParam).Return(
		&activities.ValidateAccountActivityResult{
			Valid: false,
		},
		nil,
	)

	s.env.ExecuteWorkflow(workflows.PaymentWorkFlowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(
		s.env.GetWorkflowError(),
		"account validation failed",
	)
}

func (s *UnitTestSuite) TestPaymentWorkFlowDefinition_DebitActivityFailed() {
	param := workflows.PaymentWorkFlowParam{
		OrderID:   "order-1",
		AccountID: "account-1",
		Amount:    100,
	}
	// Mock the activity with expected behavior
	s.env.OnActivity(
		activities.ValidateAccountActivityName,
		activities.ValidateAccountActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.ValidateAccountActivityResult{
			Valid: true,
		},
		nil,
	)
	debitErr := errors.New("debit activity failed")
	s.env.OnActivity(
		activities.DebitActivityName,
		activities.DebitActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		}).Return(
		nil,
		debitErr,
	)
	s.env.ExecuteWorkflow(workflows.PaymentWorkFlowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(
		s.env.GetWorkflowError(),
		debitErr.Error(),
	)
	// Credit should NOT be called when debit fails - there's nothing to compensate
	s.env.AssertActivityNotCalled(s.T(), activities.CreditActivityName)
}

func (s *UnitTestSuite) TestPaymentWorkFlowDefinition_FraudCheckTimeout() {
	param := workflows.PaymentWorkFlowParam{
		OrderID:   "order-1",
		AccountID: "account-1",
		Amount:    100,
	}
	// Mock the activity with expected behavior
	s.env.OnActivity(
		activities.ValidateAccountActivityName,
		activities.ValidateAccountActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.ValidateAccountActivityResult{
			Valid: true,
		},
		nil,
	)
	s.env.OnActivity(
		activities.DebitActivityName,
		activities.DebitActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		}).Return(
		&activities.DebitActivityResult{
			Balance: 100,
		},
		nil,
	)
	s.env.OnWorkflow(
		workflows.FraudCheckWorkflowDefinition,
		mock.Anything,
		workflows.FraudCheckWorkflowParam{
			OrderID: param.OrderID,
		},
	).Return(
		workflows.FraudCheckWorkflowResult{IsValid: true},
		nil,
	).After(11 * time.Minute)

	s.env.ExecuteWorkflow(workflows.PaymentWorkFlowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(s.env.GetWorkflowError(), "fraud check timed out")
}

func (s *UnitTestSuite) TestPaymentWorkFlowDefinition_FraudCheckInvalid() {
	param := workflows.PaymentWorkFlowParam{
		OrderID:   "order-1",
		AccountID: "account-1",
		Amount:    100,
	}
	// Mock the activity with expected behavior
	s.env.OnActivity(
		activities.ValidateAccountActivityName,
		activities.ValidateAccountActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.ValidateAccountActivityResult{
			Valid: true,
		},
		nil,
	)
	s.env.OnActivity(
		activities.DebitActivityName,
		activities.DebitActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		}).Return(
		&activities.DebitActivityResult{
			Balance: 100,
		},
		nil,
	)
	s.env.OnActivity(
		activities.CreditActivityName,
		activities.CreditActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.CreditActivityResult{
			Balance: 200,
		},
		nil,
	)
	s.env.OnWorkflow(
		workflows.FraudCheckWorkflowDefinition,
		mock.Anything,
		workflows.FraudCheckWorkflowParam{
			OrderID: param.OrderID,
		},
	).Return(
		workflows.FraudCheckWorkflowResult{IsValid: false},
		nil,
	)
	s.env.ExecuteWorkflow(workflows.PaymentWorkFlowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(s.env.GetWorkflowError(), "fraud check failed")
	s.env.AssertActivityCalled(s.T(), activities.CreditActivityName, activities.CreditActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	})
}

func (s *UnitTestSuite) TestPaymentWorkFlowDefinition_NotificationFailed() {
	param := workflows.PaymentWorkFlowParam{
		OrderID:   "order-1",
		AccountID: "account-1",
		Amount:    100,
	}
	// Mock the activity with expected behavior
	s.env.OnActivity(
		activities.ValidateAccountActivityName,
		activities.ValidateAccountActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.ValidateAccountActivityResult{
			Valid: true,
		},
		nil,
	)
	s.env.OnActivity(
		activities.DebitActivityName,
		activities.DebitActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		}).Return(
		&activities.DebitActivityResult{
			Balance: 100,
		},
		nil,
	)
	s.env.OnActivity(
		activities.CreditActivityName,
		activities.CreditActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		&activities.CreditActivityResult{
			Balance: 200,
		},
		nil,
	)
	s.env.OnWorkflow(
		workflows.FraudCheckWorkflowDefinition,
		mock.Anything,
		workflows.FraudCheckWorkflowParam{
			OrderID: param.OrderID,
		},
	).Return(
		workflows.FraudCheckWorkflowResult{IsValid: true},
		nil,
	)
	s.env.OnActivity(
		activities.NotifyPaymentActivityName,
		activities.NotifyPaymentActivityParam{
			AccountID: param.AccountID,
			Amount:    param.Amount,
		},
	).Return(
		errors.New("failed to send notification"),
	)
	s.env.ExecuteWorkflow(workflows.PaymentWorkFlowDefinition, param)
	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(s.env.GetWorkflowError(), "failed to send notification")
	s.env.AssertActivityCalled(s.T(), activities.CreditActivityName, activities.CreditActivityParam{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	})
}
