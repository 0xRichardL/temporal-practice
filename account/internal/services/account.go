package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/0xRichardL/temporal-practice/account/internal/dtos"
	"github.com/0xRichardL/temporal-practice/account/internal/models"
	"github.com/0xRichardL/temporal-practice/account/internal/seeds"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) Seeds(ctx context.Context) error {
	s.db.WithContext(ctx).Exec(fmt.Sprintf("TRUNCATE TABLE %s", models.Account{}.TableName()))
	if err := gorm.G[models.Account](s.db).CreateInBatches(ctx, &seeds.ACCOUNTS, len(seeds.ACCOUNTS)); err != nil {
		return err
	}
	return nil
}

func (s *AccountService) ValidateAccount(ctx context.Context, dto dtos.ValidateAccountRequest) (*dtos.ValidateAccountResponse, error) {
	bln, err := gorm.G[models.Account](s.db).Where("id = ?", dto.AccountID).First(ctx)
	if err != nil {
		return nil, err
	}
	return &dtos.ValidateAccountResponse{Valid: bln.Balance >= dto.Amount}, nil
}

func (s *AccountService) Debit(ctx context.Context, dto dtos.DebitRequest) (*dtos.DebitResponse, error) {
	var acc models.Account
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ?", dto.AccountID).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&acc).Error
		if err != nil {
			return err
		}
		if acc.Balance < dto.Amount {
			return fmt.Errorf("insufficient funds")
		}
		if err := tx.Model(&acc).Update("balance", acc.Balance-dto.Amount).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &dtos.DebitResponse{
		Balance: acc.Balance,
		Reason:  dto.Reason,
	}, nil
}

func (s *AccountService) Credit(ctx context.Context, dto dtos.CreditRequest) (*dtos.CreditResponse, error) {
	var acc models.Account
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ?", dto.AccountID).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&acc).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&acc).Update("balance", acc.Balance+dto.Amount).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &dtos.CreditResponse{
		Balance: acc.Balance,
		Reason:  dto.Reason,
	}, nil
}
