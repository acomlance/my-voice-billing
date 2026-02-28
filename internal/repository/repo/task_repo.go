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
	Create(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id int64) error
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

func (s *TaskRepo) Create(ctx context.Context, t *models.Task) error {
	query := fmt.Sprintf("INSERT INTO %s (token, account_id, reserved_tokens, date_create) VALUES ($1, $2, $3, NOW()) RETURNING id, date_create", db.TableTasks)
	row := s.writer.QueryRowxContext(ctx, query, t.Token, t.AccountId, t.ReservedTokens)
	err := row.Scan(&t.Id, &t.DateCreate)
	if err != nil {
		if db.IsDuplicateError(err) {
			return fmt.Errorf("task: %w", domain.ErrConflict)
		}
		if db.IsForeignKeyError(err) {
			return fmt.Errorf("task: %w", domain.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *TaskRepo) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", db.TableTasks)
	res, err := s.writer.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("task: %w", domain.ErrNotFound)
	}
	return nil
}
