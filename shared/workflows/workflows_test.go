package workflows_test

import (
	"testing"

	"github.com/0xRichardL/temporal-practice/shared/activities"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	// Workflow specific setup
	s.SetupPaymentWorkflowTest()
}

func (s *UnitTestSuite) SetupPaymentWorkflowTest() {
	s.env.RegisterActivityWithOptions(
		func(activities.ValidateAccountActivityParam) (*activities.ValidateAccountActivityResult, error) {
			return nil, nil
		},
		activity.RegisterOptions{Name: activities.ValidateAccountActivityName},
	)
	s.env.RegisterActivityWithOptions(
		func(activities.DebitActivityParam) (*activities.DebitActivityResult, error) {
			return nil, nil
		},
		activity.RegisterOptions{Name: activities.DebitActivityName},
	)
	s.env.RegisterActivityWithOptions(
		func(activities.CreditActivityParam) (*activities.CreditActivityResult, error) {
			return nil, nil
		},
		activity.RegisterOptions{Name: activities.CreditActivityName},
	)
	s.env.RegisterActivityWithOptions(
		func(activities.NotifyPaymentActivityParam) error {
			return nil
		},
		activity.RegisterOptions{Name: activities.NotifyPaymentActivityName},
	)
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}
func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
