package repository

import (
	"context"
	"testing"
	"wallet/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *DB {
	pool, err := pgxpool.New(context.Background(), "postgres://user:password@localhost:5433/walletdb_test?sslmode=disable")
	require.NoError(t, err)

	return &DB{Pool: pool}
}

func TestWalletRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewWalletRepository(db)

	t.Run("Deposit should increase balance", func(t *testing.T) {
		walletID := "test-wallet-1"
		ctx := context.Background()

		err := repo.InsertTransaction(ctx, model.TransactionRequest{
			WalletID:      walletID,
			OperationType: model.Deposit,
			Amount:        1000,
		})
		assert.NoError(t, err)

		balance, err := repo.CalculateBalance(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1000), balance)
	})

	t.Run("Withdraw should decrease balance", func(t *testing.T) {
		walletID := "test-wallet-2"
		ctx := context.Background()

		// Deposit first
		err := repo.InsertTransaction(ctx, model.TransactionRequest{
			WalletID:      walletID,
			OperationType: model.Deposit,
			Amount:        2000,
		})
		assert.NoError(t, err)

		// Then withdraw
		err = repo.InsertTransaction(ctx, model.TransactionRequest{
			WalletID:      walletID,
			OperationType: model.Withdraw,
			Amount:        500,
		})
		assert.NoError(t, err)

		balance, err := repo.CalculateBalance(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1500), balance)
	})

	t.Run("Withdraw with insufficient funds should fail", func(t *testing.T) {
		walletID := "test-wallet-3"
		ctx := context.Background()

		err := repo.InsertTransaction(ctx, model.TransactionRequest{
			WalletID:      walletID,
			OperationType: model.Withdraw,
			Amount:        1000,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient funds")
	})
}
