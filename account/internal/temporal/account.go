package temporal

import (
	"context"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"

	"github.com/0xRichardL/temporal-practice/account/internal/dtos"
	"github.com/0xRichardL/temporal-practice/account/internal/services"
	"github.com/0xRichardL/temporal-practice/shared/activities"
)

type AccountActivities struct {
	accountService *services.AccountService
}

func NewAccountActivities(accountService *services.AccountService) activities.AccountActivities {
	return &AccountActivities{
		accountService: accountService,
	}
}

func (a *AccountActivities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(a.ValidateAccount, activity.RegisterOptions{Name: activities.ValidateAccountActivityName})
	w.RegisterActivityWithOptions(a.Debit, activity.RegisterOptions{Name: activities.DebitActivityName})
	w.RegisterActivityWithOptions(a.Credit, activity.RegisterOptions{Name: activities.CreditActivityName})
}

func (a *AccountActivities) ValidateAccount(ctx context.Context, param activities.ValidateAccountActivityParam) (*activities.ValidateAccountActivityResult, error) {
	result, err := a.accountService.ValidateAccount(ctx, dtos.ValidateAccountRequest{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &activities.ValidateAccountActivityResult{
		Valid: result.Valid,
	}, nil

}

func (a *AccountActivities) Debit(ctx context.Context, param activities.DebitActivityParam) (*activities.DebitActivityResult, error) {
	result, err := a.accountService.Debit(ctx, dtos.DebitRequest{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &activities.DebitActivityResult{
		Balance: result.Balance,
	}, nil
}

func (a *AccountActivities) Credit(ctx context.Context, param activities.CreditActivityParam) (*activities.CreditActivityResult, error) {
	result, err := a.accountService.Credit(ctx, dtos.CreditRequest{
		AccountID: param.AccountID,
		Amount:    param.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &activities.CreditActivityResult{
		Balance: result.Balance,
	}, nil
}
