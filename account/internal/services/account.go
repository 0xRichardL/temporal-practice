package services

import (
	"context"

	"github.com/0xRichardL/temporal-practice/account/internal/dtos"
	"github.com/0xRichardL/temporal-practice/account/internal/models"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) ValidateAccount(ctx context.Context, dto dtos.ValidateAccountRequest) (*dtos.ValidateAccountResponse, error) {
	bln, err := gorm.G[models.Account](s.db).Where("account_id = ?", dto.AccountID).First(ctx)
	if err != nil {
		return nil, err
	}
	return &dtos.ValidateAccountResponse{Valid: bln.Balance >= dto.Amount}, nil
}

func (s *AccountService) Debit(ctx context.Context, dto dtos.DebitRequest) (*dtos.DebitResponse, error) {
	var acc models.Account
	if err := s.db.WithContext(ctx).Where("account_id = ?", dto.AccountID).First(&acc).Error; err != nil {
		return nil, err
	}
	if acc.Balance < dto.Amount {
		return nil, gorm.ErrInvalidData
	}
	if err := s.db.Model(&acc).Update("balance", acc.Balance-dto.Amount).Error; err != nil {
		return nil, err
	}
	return &dtos.DebitResponse{Balance: acc.Balance - dto.Amount}, nil
}

func (s *AccountService) Credit(ctx context.Context, dto dtos.CreditRequest) (*dtos.CreditResponse, error) {
	var acc models.Account
	if err := s.db.WithContext(ctx).Where("account_id = ?", dto.AccountID).First(&acc).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&acc).Update("balance", acc.Balance+dto.Amount).Error; err != nil {
		return nil, err
	}
	return &dtos.CreditResponse{Balance: acc.Balance + dto.Amount}, nil
}
