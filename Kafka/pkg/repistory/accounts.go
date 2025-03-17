package repistory

import (
	"context"
	"go_kafka/modul.go"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AccountRepistory struct {
	db *pgxpool.Pool
}

func NewAccountRepistory(db *pgxpool.Pool) *AccountRepistory {
	return &AccountRepistory{
		db: db,
	}
}

func (a *AccountRepistory) CraeteAccount(context context.Context, user modul.Account) (int, error) {
	var id int
	query := "INSERT INTO accounts(balance, currency, is_locked, created_at) VALUES($1, $2, $3, $4) RETURNING id"

	row := a.db.QueryRow(context, query, user.Balance, user.Currency, user.IsLocked, time.Now())

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	sql := "INSERT INTO transactions(account_id, amount, transaction_type, created_at) VALUES($1, $2, $3, $4)"
	_, err := a.db.Exec(context, sql, id, user.Balance, "deposit", time.Now())
	return id, err
}

func (a *AccountRepistory) DeleteAccount(ctx context.Context, account_id int) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return context.Background().Err()
	}
	var date time.Time
	query := "UPDATE accounts SET deleted_at=$1 WHERE id = $2 RETURNING deleted_at"
	row := tx.QueryRow(ctx, query, time.Now(), account_id)
	if err = row.Scan(&date); err != nil {
		return err
	}
	sql := "UPDATE transactions SET deleted_at=$1 WHERE account_id = $2"
	_, err = tx.Exec(ctx, sql, date, account_id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (a *AccountRepistory) GetAccountById(ctx context.Context, account_id int) (modul.Account, error) {
	query := "SELECT id, balance, currency, is_locked, created_at, deleted_at FROM accounts WHERE id = $1"
	var account modul.Account
	row := a.db.QueryRow(ctx, query, account_id)
	if err := row.Scan(&account.Id, &account.Balance,
		&account.Currency,
		&account.IsLocked, &account.CreatedAt,
		&account.DeletedAt); err != nil {
		return modul.Account{}, err
	}
	return account, nil

}

func (a *AccountRepistory) GetAllAccounts(ctx context.Context) ([]modul.Account, error) {
	query := "SELECT id, balance, currency, is_locked, created_at, deleted_at FROM accounts"
	var accounts []modul.Account
	rows, err := a.db.Query(ctx, query)
	if err != nil {
		return []modul.Account{}, err
	}
	for rows.Next() {
		var account modul.Account
		if err := rows.Scan(&account.Id, &account.Balance,
			&account.Currency,
			&account.IsLocked, &account.CreatedAt,
			&account.DeletedAt); err != nil {
			return []modul.Account{}, err
		}
		if account.DeletedAt == nil {
			accounts = append(accounts, account)
		}
	}
	return accounts, nil

}
