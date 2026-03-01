package logic

import (
	"context"
	"errors"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"

	"github.com/google/uuid"
)

type TaskLogic struct {
	taskRepo    repo.TaskRepository
	accountRepo repo.AccountRepository
}

func NewTaskLogic(taskRepo repo.TaskRepository, accountRepo repo.AccountRepository) *TaskLogic {
	return &TaskLogic{taskRepo: taskRepo, accountRepo: accountRepo}
}

func (l *TaskLogic) Create(ctx context.Context, t *models.Task) error {
	if t.ReservedTokens <= 0 {
		return domain.ErrInvalid
	}
	if _, err := l.accountRepo.GetByID(ctx, t.AccountId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	t.Token = uuid.New().String()
	return l.taskRepo.CreateWithReserveUpdate(ctx, t)
}

func (l *TaskLogic) Delete(ctx context.Context, token string, closedTokens int64) error {
	if closedTokens < 0 {
		return domain.ErrInvalid
	}
	return l.taskRepo.DeleteByTokenWithReserveUpdate(ctx, token, closedTokens)
}
