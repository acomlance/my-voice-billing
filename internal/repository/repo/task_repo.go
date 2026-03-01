package repo

import (
	"context"
	"fmt"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/db"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	GetByToken(ctx context.Context, token string) (*models.Task, error)
	ListByAccountID(ctx context.Context, accountID int64) ([]models.Task, error)
	CreateWithReserveUpdate(ctx context.Context, t *models.Task) error
	DeleteByTokenWithReserveUpdate(ctx context.Context, token string, closedTokens int64) error
}

type TaskRepo struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewTaskRepo(m *db.Manager) *TaskRepo {
	return &TaskRepo{reader: m.Reader(), writer: m.Writer()}
}

func (s *TaskRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := fmt.Sprintf("SELECT id, token, account_id, reserved_tokens, date_create FROM %s WHERE id = $1", db.TableTasks)
	var t models.Task
	err := s.reader.GetContext(ctx, &t, query, id)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("task: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &t, nil
}

func (s *TaskRepo) GetByToken(ctx context.Context, token string) (*models.Task, error) {
	query := fmt.Sprintf("SELECT id, token, account_id, reserved_tokens, date_create FROM %s WHERE token = $1", db.TableTasks)
	var t models.Task
	err := s.reader.GetContext(ctx, &t, query, token)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("task: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &t, nil
}

func (s *TaskRepo) ListByAccountID(ctx context.Context, accountID int64) ([]models.Task, error) {
	query := fmt.Sprintf("SELECT id, token, account_id, reserved_tokens, date_create FROM %s WHERE account_id = $1", db.TableTasks)
	var list []models.Task
	err := s.reader.SelectContext(ctx, &list, query, accountID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TaskRepo) CreateWithReserveUpdate(ctx context.Context, t *models.Task) error {
	tx, err := s.writer.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	lock := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 FOR UPDATE", db.TableAccounts)
	if _, err := tx.ExecContext(ctx, lock, t.AccountId); err != nil {
		return err
	}
	var sum int64
	qSum := fmt.Sprintf("SELECT COALESCE(SUM(reserved_tokens), 0) FROM %s WHERE account_id = $1", db.TableTasks)
	if err := tx.GetContext(ctx, &sum, qSum, t.AccountId); err != nil {
		return err
	}
	newReserve := sum + t.ReservedTokens

	upd := fmt.Sprintf("UPDATE %s SET reserve = $2, date_update = NOW() WHERE id = $1", db.TableAccounts)
	_, err = tx.ExecContext(ctx, upd, t.AccountId, newReserve)
	if err != nil {
		if db.IsCheckConstraintError(err) {
			return fmt.Errorf("account: %w", domain.ErrInsufficientBalance)
		}
		return err
	}

	ins := fmt.Sprintf("INSERT INTO %s (token, account_id, reserved_tokens, date_create) VALUES ($1, $2, $3, NOW()) RETURNING id, date_create", db.TableTasks)
	row := tx.QueryRowxContext(ctx, ins, t.Token, t.AccountId, t.ReservedTokens)
	if err := row.Scan(&t.Id, &t.DateCreate); err != nil {
		if db.IsDuplicateError(err) {
			return fmt.Errorf("task: %w", domain.ErrConflict)
		}
		return err
	}
	return tx.Commit()
}

func (s *TaskRepo) DeleteByTokenWithReserveUpdate(ctx context.Context, token string, closedTokens int64) error {
	tx, err := s.writer.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var t models.Task
	qGet := fmt.Sprintf("SELECT id, token, account_id, reserved_tokens FROM %s WHERE token = $1", db.TableTasks)
	if err := tx.GetContext(ctx, &t, qGet, token); err != nil {
		if db.IsNoRowsError(err) {
			return fmt.Errorf("task: %w", domain.ErrNotFound)
		}
		return err
	}

	upd := fmt.Sprintf("UPDATE %s SET reserve = reserve - $2, balance = balance - $3, date_update = NOW() WHERE id = $1", db.TableAccounts)
	_, err = tx.ExecContext(ctx, upd, t.AccountId, t.ReservedTokens, closedTokens)
	if err != nil {
		if db.IsCheckConstraintError(err) {
			return fmt.Errorf("account: %w", domain.ErrInsufficientBalance)
		}
		return err
	}

	del := fmt.Sprintf("DELETE FROM %s WHERE id = $1", db.TableTasks)
	if _, err := tx.ExecContext(ctx, del, t.Id); err != nil {
		return err
	}
	return tx.Commit()
}
