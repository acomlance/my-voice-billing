package repo

import (
	"context"
	"fmt"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/db"

	"github.com/jmoiron/sqlx"
)

type TokenRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Token, error)
	GetByToken(ctx context.Context, token string) (*models.Token, error)
	List(ctx context.Context) ([]models.Token, error)
	Create(ctx context.Context, t *models.Token) error
	Update(ctx context.Context, t *models.Token) error
	Delete(ctx context.Context, id int64) error
}

type TokenRepo struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewTokenRepo(m *db.Manager) *TokenRepo {
	return &TokenRepo{reader: m.Reader(), writer: m.Writer()}
}

func (s *TokenRepo) GetByID(ctx context.Context, id int64) (*models.Token, error) {
	query := fmt.Sprintf("SELECT id, user_id, token, date_create, date_update FROM %s WHERE id = $1", db.TableToken)
	var t models.Token
	err := s.reader.GetContext(ctx, &t, query, id)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("token: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &t, nil
}

func (s *TokenRepo) GetByToken(ctx context.Context, token string) (*models.Token, error) {
	query := fmt.Sprintf("SELECT id, user_id, token, date_create, date_update FROM %s WHERE token = $1", db.TableToken)
	var t models.Token
	err := s.reader.GetContext(ctx, &t, query, token)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("token: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &t, nil
}

func (s *TokenRepo) List(ctx context.Context) ([]models.Token, error) {
	query := fmt.Sprintf("SELECT id, user_id, token, date_create, date_update FROM %s", db.TableToken)
	var list []models.Token
	err := s.reader.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *TokenRepo) Create(ctx context.Context, t *models.Token) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token, date_create, date_update) VALUES ($1, $2, NOW(), NOW()) RETURNING id, date_create, date_update", db.TableToken)
	row := s.writer.QueryRowxContext(ctx, query, t.UserId, t.Token)
	err := row.Scan(&t.Id, &t.DateCreate, &t.DateUpdate)
	if err != nil {
		if db.IsDuplicateError(err) {
			return fmt.Errorf("token: %w", domain.ErrConflict)
		}
		return err
	}
	return nil
}

func (s *TokenRepo) Update(ctx context.Context, t *models.Token) error {
	query := fmt.Sprintf("UPDATE %s SET token = $2, date_update = NOW() WHERE id = $1 RETURNING date_update", db.TableToken)
	row := s.writer.QueryRowxContext(ctx, query, t.Id, t.Token)
	err := row.Scan(&t.DateUpdate)
	if err != nil {
		if db.IsNoRowsError(err) {
			return fmt.Errorf("token: %w", domain.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *TokenRepo) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", db.TableToken)
	res, err := s.writer.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("token: %w", domain.ErrNotFound)
	}
	return nil
}
