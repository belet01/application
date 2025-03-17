package repistory

import (
	"context"
	"fmt"

	"go_kafka/modul.go"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionPostgres struct {
	db *pgxpool.Pool
}

func NewTransactionPostgres(db *pgxpool.Pool) *TransactionPostgres {
	return &TransactionPostgres{
		db: db,
	}
}

func (t *TransactionPostgres) Deposit(ctx context.Context, balance float64, account_id int) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	query := "UPDATE accounts SET balance=balance + $1 WHERE id =$2"
	_, err = tx.Exec(ctx, query, balance, account_id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	sql := "INSERT INTO transactions(account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'deposit', NOW())"
	_, err = tx.Exec(ctx, sql, account_id, balance)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (t *TransactionPostgres) Withdraw(ctx context.Context, amount float64, accountID int) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	var currentBalance float32
	query := "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE"
	row := tx.QueryRow(ctx, query, accountID)
	if err = row.Scan(&currentBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("bakiye alınamadı: %w", err)
	}

	if currentBalance < float32(amount) {
		tx.Rollback(ctx)
		return fmt.Errorf("yetersiz bakiye! Mevcut bakiye: %.2f", currentBalance)
	}

	updateQuery := "UPDATE accounts SET balance = balance - $1 WHERE id = $2"
	_, err = tx.Exec(ctx, updateQuery, amount, accountID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("bakiye güncellenemedi: %w", err)
	}
	insertQuery := "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'withdraw', NOW())"
	_, err = tx.Exec(ctx, insertQuery, accountID, amount)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("işlem kaydı eklenemedi: %w", err)
	}
	return tx.Commit(ctx)
}

func (t *TransactionPostgres) Transfer(ctx context.Context, param_id int, amount float64, accountID int) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	var currentBalance float32
	query := "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE"
	row := tx.QueryRow(ctx, query, accountID)
	if err = row.Scan(&currentBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("bakiye alınamadı: %w", err)
	}

	if currentBalance < float32(amount) {
		tx.Rollback(ctx)
		return fmt.Errorf("yetersiz bakiye! Mevcut bakiye: %.2f", currentBalance)
	}

	updateQuery := "UPDATE accounts SET balance = balance - $1 WHERE id = $2"
	_, err = tx.Exec(ctx, updateQuery, amount, accountID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("bakiye güncellenemedi: %w", err)
	}
	insertQuery := "INSERT INTO transactions (account_id, amount, transaction_type, created_at) VALUES ($1, $2, 'transfer', NOW())"
	_, err = tx.Exec(ctx, insertQuery, accountID, amount)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("işlem kaydı eklenemedi: %w", err)
	}
	return tx.Commit(ctx)
}

func (t *TransactionPostgres) GetTransactions(ctx context.Context, account_id int) (modul.TransactionResponse, error) {
	query := "SELECT id, amount, transaction_type, created_at FROM transactions WHERE account_id =$1"
	rows, err := t.db.Query(ctx, query, account_id)
	if err != nil {
		return modul.TransactionResponse{}, err
	}
	var accounts []modul.Transactions
	for rows.Next() {
		var account modul.Transactions
		err := rows.Scan(&account.Id,
			&account.Amount,
			&account.TransactionType,
			&account.CreatedAt)
		if err != nil {
			return modul.TransactionResponse{}, err
		}
		accounts = append(accounts, account)
	}
	accountget := modul.TransactionResponse{
		AccountId:   account_id,
		Transaction: accounts,
	}
	return accountget, nil

}
