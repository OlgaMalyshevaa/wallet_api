// package service

// import (
// 	"context"
// 	"errors"
// 	"wallet/internal/model"
// 	"wallet/internal/repository"
// )

// type WalletService struct {
// 	repo *repository.WalletRepository
// }

// func NewWalletService(repo *repository.WalletRepository) *WalletService {
// 	return &WalletService{repo: repo}
// }

// func (s *WalletService) PerformTransaction(ctx context.Context, req model.TransactionRequest) error {
// 	if req.Amount <= 0 {
// 		return errors.New("amount must be positive")
// 	}
// 	return s.repo.InsertTransaction(ctx, req)
// }

// func (s *WalletService) GetBalance(ctx context.Context, id string) (int64, error) {
// 	return s.repo.CalculateBalance(ctx, id)
// }

package service

import (
	"context"
	"errors"
	"wallet/internal/model"
	"wallet/internal/repository"
)

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(db *repository.DB) *WalletService {
	return &WalletService{repo: repository.NewWalletRepository(db)}
}

func (s *WalletService) PerformTransaction(ctx context.Context, req model.TransactionRequest) error {
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.repo.InsertTransaction(ctx, req)
}

func (s *WalletService) GetBalance(ctx context.Context, id string) (int64, error) {
	return s.repo.CalculateBalance(ctx, id)
}
