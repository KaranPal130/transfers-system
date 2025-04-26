package service

import (
	"context"
	"errors"

	"github.com/KaranPal130/transfers-system/internal/models"
	repository "github.com/KaranPal130/transfers-system/internal/repositories"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidInitialBalance = errors.New("invalid initial balance")
	ErrAccountAlreadyExists  = errors.New("account already exists")
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req models.AccountCreateRequest) error {
	initialBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return ErrInvalidInitialBalance
	}

	if initialBalance.LessThan(decimal.Zero) {
		return ErrInvalidInitialBalance
	}

	_, err = s.accountRepo.GetByID(ctx, req.AccountID)
	if err == nil {
		return ErrAccountAlreadyExists
	} else if !errors.Is(err, repository.ErrAccountNotFound) {
		return err
	}

	account := models.Account{
		AccountID: req.AccountID,
		Balance:   initialBalance,
	}

	return s.accountRepo.Create(ctx, account)
}

func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (models.Account, error) {
	return s.accountRepo.GetByID(ctx, accountID)
}
