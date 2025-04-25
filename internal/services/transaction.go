package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KaranPal130/transfers-system/internal/models"
	repository "github.com/KaranPal130/transfers-system/internal/repositories"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSameSourceAndDest   = errors.New("source and destination accounts must be different")
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	db              *sql.DB
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	db *sql.DB,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// CreateTransaction processes a transfer between accounts
func (s *TransactionService) CreateTransaction(ctx context.Context, req models.TransactionRequest) error {
	// Validate request
	if req.SourceAccountID == req.DestinationAccountID {
		return ErrSameSourceAndDest
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return ErrInvalidAmount
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	// Start database transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Get source account (with locking to prevent race conditions)
	sourceAccount, err := s.accountRepo.GetByID(ctx, req.SourceAccountID)
	if err != nil {
		return err
	}

	// Get destination account
	destAccount, err := s.accountRepo.GetByID(ctx, req.DestinationAccountID)
	if err != nil {
		return err
	}

	// Check if source has enough balance
	if sourceAccount.Balance.LessThan(amount) {
		return ErrInsufficientBalance
	}

	// Update source account balance
	newSourceBalance := sourceAccount.Balance.Sub(amount)
	err = s.accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance)
	if err != nil {
		return err
	}

	// Update destination account balance
	newDestBalance := destAccount.Balance.Add(amount)
	err = s.accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance)
	if err != nil {
		return err
	}

	// Record transaction
	transaction := models.Transaction{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
	}

	err = s.transactionRepo.Create(ctx, tx, transaction)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}
