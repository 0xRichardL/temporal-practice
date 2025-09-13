package services

import (
	"github.com/0xRichardL/temporal-practice/account/internal/assets"
	"github.com/0xRichardL/temporal-practice/account/internal/models"
)

type AccountService struct {
	db []models.Balance
}

func NewAccountService() *AccountService {
	return &AccountService{
		db: assets.BALANCES,
	}
}

func (s *AccountService) Debit() {}

func (s *AccountService) Credit() {}

func (s *AccountService) Rollback() {}
