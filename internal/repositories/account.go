package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KaranPal130/transfers-system/internal/models"
	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound = errors.New("Account not Found")
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Create(ctx context.Context, account models.Account) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, account.AccountID, account.Balance)
	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID int64) (models.Account, error) {
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1`

	var account models.Account
	var balanceStr string

	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&account.AccountID, &balanceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, ErrAccountNotFound
		}
		return models.Account{}, err
	}

	account.Balance, err = decimal.NewFromString(balanceStr)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, accountID int64) (models.Account, error) {
    query := `SELECT account_id, balance FROM accounts WHERE account_id = $1 FOR UPDATE`
    var account models.Account
    var balanceStr string
    err := tx.QueryRowContext(ctx, query, accountID).Scan(&account.AccountID, &balanceStr)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Account{}, ErrAccountNotFound
        }
        return models.Account{}, err
    }
    account.Balance, err = decimal.NewFromString(balanceStr)
    if err != nil {
        return models.Account{}, err
    }
    return account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID int64, newBalance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1 WHERE account_id = $2`
	result, err := tx.ExecContext(ctx, query, newBalance.String(), accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAccountNotFound
	}

	return nil
}