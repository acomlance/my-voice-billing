package logic

import (
	"context"

	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"
)

type TokenLogic struct {
	repo repo.TokenRepository
}

func NewTokenLogic(repo repo.TokenRepository) *TokenLogic {
	return &TokenLogic{repo: repo}
}

func (l *TokenLogic) GetByID(ctx context.Context, id int64) (*models.Token, error) {
	return l.repo.GetByID(ctx, id)
}

func (l *TokenLogic) GetByToken(ctx context.Context, token string) (*models.Token, error) {
	return l.repo.GetByToken(ctx, token)
}

func (l *TokenLogic) List(ctx context.Context) ([]models.Token, error) {
	return l.repo.List(ctx)
}

func (l *TokenLogic) Create(ctx context.Context, t *models.Token) error {
	return l.repo.Create(ctx, t)
}

func (l *TokenLogic) Update(ctx context.Context, t *models.Token) error {
	return l.repo.Update(ctx, t)
}

func (l *TokenLogic) Delete(ctx context.Context, id int64) error {
	return l.repo.Delete(ctx, id)
}
