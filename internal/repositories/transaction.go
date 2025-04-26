package repository

import (
	"context"
	"database/sql"

	"github.com/KaranPal130/transfers-system/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction models.Transaction) error {
	query := `
		INSERT INTO transactions (source_account_id, destination_account_id, amount)
		VALUES ($1, $2, $3)
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		transaction.SourceAccountID,
		transaction.DestinationAccountID,
		transaction.Amount,
	)

	return err
}
