package logic

import (
	"context"
	"errors"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"
)

type TaskLogic struct {
	taskRepo  repo.TaskRepository
	tokenRepo repo.TokenRepository
}

func NewTaskLogic(taskRepo repo.TaskRepository, tokenRepo repo.TokenRepository) *TaskLogic {
	return &TaskLogic{taskRepo: taskRepo, tokenRepo: tokenRepo}
}

func (l *TaskLogic) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	return l.taskRepo.GetByID(ctx, id)
}

func (l *TaskLogic) ListByTokenID(ctx context.Context, tokenID int64) ([]models.Task, error) {
	return l.taskRepo.ListByTokenID(ctx, tokenID)
}

func (l *TaskLogic) Create(ctx context.Context, t *models.Task) error {
	if _, err := l.tokenRepo.GetByID(ctx, t.TokenId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return l.taskRepo.Create(ctx, t)
}

func (l *TaskLogic) Update(ctx context.Context, t *models.Task) error {
	return l.taskRepo.Update(ctx, t)
}

func (l *TaskLogic) Delete(ctx context.Context, id int64) error {
	return l.taskRepo.Delete(ctx, id)
}
