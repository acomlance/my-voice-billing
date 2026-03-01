package logic

import (
	"context"
	"errors"

	"my-voice-billing/internal/domain"
	"my-voice-billing/internal/models"
	"my-voice-billing/internal/repository/repo"
)

type TransactionLogic struct {
	transactionRepo repo.TransactionRepository
	accountRepo     repo.AccountRepository
}

func NewTransactionLogic(transactionRepo repo.TransactionRepository, accountRepo repo.AccountRepository) *TransactionLogic {
	return &TransactionLogic{transactionRepo: transactionRepo, accountRepo: accountRepo}
}

func (l *TransactionLogic) Create(ctx context.Context, t *models.Transaction) error {
	if _, err := l.accountRepo.GetByID(ctx, t.AccountId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return l.transactionRepo.Create(ctx, t)
}

func (l *TransactionLogic) GetByID(ctx context.Context, id int64) (*models.Transaction, error) {
	return l.transactionRepo.GetByID(ctx, id)
}

func (l *TransactionLogic) List(ctx context.Context, accountID int64, limit, offset int) ([]models.Transaction, error) {
	return l.transactionRepo.List(ctx, accountID, limit, offset)
}

func (l *TransactionLogic) Update(ctx context.Context, t *models.Transaction) error {
	return l.transactionRepo.Update(ctx, t)
}
