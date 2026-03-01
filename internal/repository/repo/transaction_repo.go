package repo

import (
	"context"
	"fmt"

	"my-voice-billing/internal/constants"
	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/db"

	"github.com/jmoiron/sqlx"
)

const txCols = "id, account_id, status, amount, tokens, tokens_after, payment_type, payment_method, payment_data, description, date_create, date_update"

type TransactionRepository interface {
	Create(ctx context.Context, t *models.Transaction) error
	GetByID(ctx context.Context, id int64) (*models.Transaction, error)
	List(ctx context.Context, accountID int64, limit, offset int) ([]models.Transaction, error)
	Update(ctx context.Context, t *models.Transaction) error
}

type TransactionRepo struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewTransactionRepo(m *db.Manager) *TransactionRepo {
	return &TransactionRepo{reader: m.Reader(), writer: m.Writer()}
}

func (r *TransactionRepo) GetByID(ctx context.Context, id int64) (*models.Transaction, error) {
	q := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", txCols, db.TableTransactions)
	var t models.Transaction
	if err := r.reader.GetContext(ctx, &t, q, id); err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("transaction: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepo) Create(ctx context.Context, t *models.Transaction) error {
	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if t.Status == constants.TxStatusApproved {
		var balanceAfter int64
		upd := fmt.Sprintf("UPDATE %s SET balance = balance + $2, date_update = NOW() WHERE id = $1 RETURNING balance", db.TableAccounts)
		if err := tx.GetContext(ctx, &balanceAfter, upd, t.AccountId, t.Tokens); err != nil {
			if db.IsCheckConstraintError(err) {
				return fmt.Errorf("account: %w", domain.ErrInsufficientBalance)
			}
			return err
		}
		t.TokensAfter = &balanceAfter
	}

	ins := fmt.Sprintf("INSERT INTO %s (account_id, status, amount, tokens, tokens_after, payment_type, payment_method, payment_data, description, date_create, date_update) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id, date_create, date_update", db.TableTransactions)
	row := tx.QueryRowxContext(ctx, ins, t.AccountId, t.Status, t.Amount, t.Tokens, t.TokensAfter, t.PaymentType, t.PaymentMethod, t.PaymentData, t.Description)
	if err := row.Scan(&t.Id, &t.DateCreate, &t.DateUpdate); err != nil {
		if db.IsForeignKeyError(err) {
			return fmt.Errorf("transaction: %w", domain.ErrNotFound)
		}
		return err
	}
	return tx.Commit()
}

func (r *TransactionRepo) List(ctx context.Context, accountID int64, limit, offset int) ([]models.Transaction, error) {
	q := fmt.Sprintf("SELECT %s FROM %s WHERE account_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", txCols, db.TableTransactions)
	var list []models.Transaction
	if err := r.reader.SelectContext(ctx, &list, q, accountID, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *TransactionRepo) Update(ctx context.Context, t *models.Transaction) error {
	old, err := r.GetByID(ctx, t.Id)
	if err != nil {
		return err
	}

	tx, err := r.writer.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if old.Status != constants.TxStatusApproved && t.Status == constants.TxStatusApproved {
		upd := fmt.Sprintf("UPDATE %s SET balance = balance + $2, date_update = NOW() WHERE id = $1", db.TableAccounts)
		if _, err := tx.ExecContext(ctx, upd, old.AccountId, old.Tokens); err != nil {
			if db.IsCheckConstraintError(err) {
				return fmt.Errorf("account: %w", domain.ErrInsufficientBalance)
			}
			return err
		}
		var newBal int64
		qBal := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", db.TableAccounts)
		if err := tx.GetContext(ctx, &newBal, qBal, old.AccountId); err != nil {
			return err
		}
		t.TokensAfter = &newBal
	} else if old.Status == constants.TxStatusApproved && t.Status == constants.TxStatusCancelled {
		upd := fmt.Sprintf("UPDATE %s SET balance = balance - $2, date_update = NOW() WHERE id = $1", db.TableAccounts)
		if _, err := tx.ExecContext(ctx, upd, old.AccountId, old.Tokens); err != nil {
			if db.IsCheckConstraintError(err) {
				return fmt.Errorf("account: %w", domain.ErrInsufficientBalance)
			}
			return err
		}
		t.TokensAfter = nil
	}

	upd := fmt.Sprintf("UPDATE %s SET status = $2, payment_method = $3, payment_data = $4, description = $5, tokens_after = $6, date_update = NOW() WHERE id = $1", db.TableTransactions)
	if _, err := tx.ExecContext(ctx, upd, t.Id, t.Status, t.PaymentMethod, t.PaymentData, t.Description, t.TokensAfter); err != nil {
		return err
	}
	return tx.Commit()
}
