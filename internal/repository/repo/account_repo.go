package repo

import (
	"context"
	"fmt"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/db"

	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Account, error)
	List(ctx context.Context) ([]models.Account, error)
	Create(ctx context.Context, a *models.Account) error
	Update(ctx context.Context, a *models.Account) error
	Delete(ctx context.Context, id int64) error
}

type AccountRepo struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewAccountRepo(m *db.Manager) *AccountRepo {
	return &AccountRepo{reader: m.Reader(), writer: m.Writer()}
}

func (s *AccountRepo) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	query := fmt.Sprintf("SELECT id, balance, reserve, state, date_create, date_update FROM %s WHERE id = $1", db.TableAccounts)
	var a models.Account
	err := s.reader.GetContext(ctx, &a, query, id)
	if err != nil {
		if db.IsNoRowsError(err) {
			return nil, fmt.Errorf("account: %w", domain.ErrNotFound)
		}
		return nil, err
	}
	return &a, nil
}

func (s *AccountRepo) List(ctx context.Context) ([]models.Account, error) {
	query := fmt.Sprintf("SELECT id, balance, reserve, state, date_create, date_update FROM %s", db.TableAccounts)
	var list []models.Account
	err := s.reader.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *AccountRepo) Create(ctx context.Context, a *models.Account) error {
	query := fmt.Sprintf("INSERT INTO %s (id, balance, reserve, state, date_create, date_update) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING date_create, date_update", db.TableAccounts)
	row := s.writer.QueryRowxContext(ctx, query, a.Id, a.Balance, a.Reserve, a.State)
	err := row.Scan(&a.DateCreate, &a.DateUpdate)
	if err != nil {
		if db.IsDuplicateError(err) {
			return fmt.Errorf("account: %w", domain.ErrConflict)
		}
		return err
	}
	return nil
}

func (s *AccountRepo) Update(ctx context.Context, a *models.Account) error {
	query := fmt.Sprintf("UPDATE %s SET balance = $2, reserve = $3, state = $4, date_update = NOW() WHERE id = $1 RETURNING date_update", db.TableAccounts)
	row := s.writer.QueryRowxContext(ctx, query, a.Id, a.Balance, a.Reserve, a.State)
	err := row.Scan(&a.DateUpdate)
	if err != nil {
		if db.IsNoRowsError(err) {
			return fmt.Errorf("account: %w", domain.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *AccountRepo) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", db.TableAccounts)
	res, err := s.writer.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("account: %w", domain.ErrNotFound)
	}
	return nil
}
