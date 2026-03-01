package handlers

import (
	"context"

	"my-voice-billing/internal/models"
	"my-voice-billing/internal/transport/grpc/pb"
)

type TransactionLogic interface {
	Create(ctx context.Context, t *models.Transaction) error
	Update(ctx context.Context, t *models.Transaction) error
}

type TransactionServer struct {
	pb.UnimplementedTransactionServiceServer
	transactionLogic TransactionLogic
}

func NewTransactionServer(transactionLogic TransactionLogic) *TransactionServer {
	return &TransactionServer{transactionLogic: transactionLogic}
}

func (s *TransactionServer) Create(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	t := &models.Transaction{
		AccountId:     req.GetAccountId(),
		Status:        int16(req.GetStatus()),
		Amount:        req.GetAmount(),
		Tokens:        req.GetTokens(),
		PaymentType:   int16(req.GetPaymentType()),
		PaymentMethod: int16(req.GetPaymentMethod()),
		PaymentData:   req.GetPaymentData(),
		Description:   req.GetDescription(),
	}
	if err := s.transactionLogic.Create(ctx, t); err != nil {
		return nil, handleErr(err, "Create")
	}
	return &pb.CreateTransactionResponse{Id: t.Id}, nil
}

func (s *TransactionServer) Update(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.UpdateTransactionResponse, error) {
	t := &models.Transaction{
		Id:            req.GetId(),
		Status:        int16(req.GetStatus()),
		PaymentMethod: int16(req.GetPaymentMethod()),
		PaymentData:   req.GetPaymentData(),
		Description:   req.GetDescription(),
	}
	if err := s.transactionLogic.Update(ctx, t); err != nil {
		return nil, handleErr(err, "Update")
	}
	return &pb.UpdateTransactionResponse{}, nil
}
