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

// AccountService handles business logic for accounts
type AccountService struct {
	accountRepo *repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new account with initial balance
func (s *AccountService) CreateAccount(ctx context.Context, req models.AccountCreateRequest) error {
	// Validate initial balance
	initialBalance, err := decimal.NewFromString(req.InitialBalance)
	if err != nil {
		return ErrInvalidInitialBalance
	}

	if initialBalance.LessThan(decimal.Zero) {
		return ErrInvalidInitialBalance
	}

	// Check if account already exists
	_, err = s.accountRepo.GetByID(ctx, req.AccountID)
	if err == nil {
		return ErrAccountAlreadyExists
	} else if !errors.Is(err, repository.ErrAccountNotFound) {
		return err
	}

	// Create account
	account := models.Account{
		AccountID: req.AccountID,
		Balance:   initialBalance,
	}

	return s.accountRepo.Create(ctx, account)
}

// GetAccount retrieves account information
func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (models.Account, error) {
	return s.accountRepo.GetByID(ctx, accountID)
}
