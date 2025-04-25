package repository

import (
	"context"
	"errors"
	"wallet/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

type WalletRepository struct {
	db *DB
}

func NewWalletRepository(db *DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) InsertTransaction(ctx context.Context, req model.TransactionRequest) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO wallets(wallet_id) VALUES($1) ON CONFLICT DO NOTHING", req.WalletID)
	if err != nil {
		return err
	}

	if req.OperationType == model.Withdraw {
		balance, err := r.CalculateBalance(ctx, req.WalletID)
		if err != nil {
			return err
		}
		if balance < req.Amount {
			return errors.New("insufficient funds")
		}
	}

	_, err = tx.Exec(ctx, `INSERT INTO transactions(wallet_id, operation_type, amount) VALUES($1, $2, $3)`,
		req.WalletID, string(req.OperationType), req.Amount)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) CalculateBalance(ctx context.Context, walletID string) (int64, error) {
	var balance int64
	err := r.db.Pool.QueryRow(ctx, `
        SELECT COALESCE(SUM(CASE 
            WHEN operation_type = 'DEPOSIT' THEN amount
            WHEN operation_type = 'WITHDRAW' THEN -amount
            ELSE 0 END), 0)
        FROM transactions WHERE wallet_id = $1`, walletID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
