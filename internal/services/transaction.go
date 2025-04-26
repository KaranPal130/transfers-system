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

type TransactionService struct {
	db              *sql.DB
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
}

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

func (s *TransactionService) CreateTransaction(ctx context.Context, req models.TransactionRequest) error {
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	sourceAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, req.SourceAccountID)
	if err != nil {
		return err
	}

	destAccount, err := s.accountRepo.GetByIDForUpdate(ctx, tx, req.DestinationAccountID)
	if err != nil {
		return err
	}

	if sourceAccount.Balance.LessThan(amount) {
		return ErrInsufficientBalance
	}

	newSourceBalance := sourceAccount.Balance.Sub(amount)
	err = s.accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance)
	if err != nil {
		return err
	}

	newDestBalance := destAccount.Balance.Add(amount)
	err = s.accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance)
	if err != nil {
		return err
	}

	transaction := models.Transaction{
		SourceAccountID:      req.SourceAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
	}

	err = s.transactionRepo.Create(ctx, tx, transaction)
	if err != nil {
		return err
	}

	return tx.Commit()
}
