package logic

import (
	"context"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"
)

type AccountLogic struct {
	repo repo.AccountRepository
}

func NewAccountLogic(repo repo.AccountRepository) *AccountLogic {
	return &AccountLogic{repo: repo}
}

func (l *AccountLogic) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	return l.repo.GetByID(ctx, id)
}

func (l *AccountLogic) List(ctx context.Context) ([]models.Account, error) {
	return l.repo.List(ctx)
}

func (l *AccountLogic) Create(ctx context.Context, a *models.Account) error {
	if a.Balance < a.Reserve {
		return domain.ErrInvalid
	}
	return l.repo.Create(ctx, a)
}

func (l *AccountLogic) Update(ctx context.Context, a *models.Account) error {
	if a.Balance < a.Reserve {
		return domain.ErrInvalid
	}
	return l.repo.Update(ctx, a)
}

func (l *AccountLogic) Delete(ctx context.Context, id int64) error {
	return l.repo.Delete(ctx, id)
}
