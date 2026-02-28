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
	ListByTokenID(ctx context.Context, tokenID int64) ([]models.Task, error)
	Create(ctx context.Context, t *models.Task) error
	Update(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id int64) error
}

type TaskService struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewTaskService(m *db.Manager) *TaskService {
	return &TaskService{reader: m.Reader(), writer: m.Writer()}
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := fmt.Sprintf("SELECT id, token_id, task, date_create, date_update FROM %s WHERE id = $1", db.TableTask)
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

func (s *TaskService) ListByTokenID(ctx context.Context, tokenID int64) ([]models.Task, error) {
	query := fmt.Sprintf("SELECT id, token_id, task, date_create, date_update FROM %s WHERE token_id = $1", db.TableTask)
	var list []models.Task
	err := s.reader.SelectContext(ctx, &list, query, tokenID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TaskService) Create(ctx context.Context, t *models.Task) error {
	query := fmt.Sprintf("INSERT INTO %s (token_id, task, date_create, date_update) VALUES ($1, $2, NOW(), NOW()) RETURNING id, date_create, date_update", db.TableTask)
	row := s.writer.QueryRowxContext(ctx, query, t.TokenId, t.Task)
	err := row.Scan(&t.Id, &t.DateCreate, &t.DateUpdate)
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

func (s *TaskService) Update(ctx context.Context, t *models.Task) error {
	query := fmt.Sprintf("UPDATE %s SET task = $2, date_update = NOW() WHERE id = $1 RETURNING date_update", db.TableTask)
	row := s.writer.QueryRowxContext(ctx, query, t.Id, t.Task)
	err := row.Scan(&t.DateUpdate)
	if err != nil {
		if db.IsNoRowsError(err) {
			return fmt.Errorf("task: %w", domain.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", db.TableTask)
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
