package logic

import (
	"context"
	"errors"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"
)

type TaskLogic struct {
	taskRepo    repo.TaskRepository
	accountRepo repo.AccountRepository
}

func NewTaskLogic(taskRepo repo.TaskRepository, accountRepo repo.AccountRepository) *TaskLogic {
	return &TaskLogic{taskRepo: taskRepo, accountRepo: accountRepo}
}

func (l *TaskLogic) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	return l.taskRepo.GetByID(ctx, id)
}

func (l *TaskLogic) GetByToken(ctx context.Context, token string) (*models.Task, error) {
	return l.taskRepo.GetByToken(ctx, token)
}

func (l *TaskLogic) ListByAccountID(ctx context.Context, accountID int64) ([]models.Task, error) {
	return l.taskRepo.ListByAccountID(ctx, accountID)
}

func (l *TaskLogic) Create(ctx context.Context, t *models.Task) error {
	if _, err := l.accountRepo.GetByID(ctx, t.AccountId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return l.taskRepo.Create(ctx, t)
}

func (l *TaskLogic) Delete(ctx context.Context, id int64) error {
	return l.taskRepo.Delete(ctx, id)
}
