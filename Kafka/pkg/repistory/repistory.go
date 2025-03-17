package repistory

import (
	"context"

	"go_kafka/modul.go"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Accounts interface {
	CraeteAccount(context context.Context, user modul.Account) (int, error)
	DeleteAccount(ctx context.Context, account_id int) error
	GetAccountById(ctx context.Context, account_id int) (modul.Account, error)
	GetAllAccounts(ctx context.Context) ([]modul.Account, error)
}

type Transactions interface {
	Deposit(ctx context.Context, balance float64, account_id int) error
	CheckAccountActive(ctx context.Context, accountID int) (bool, error)
	Withdraw(ctx context.Context, amount float64, accountID int) error
	Transfer(ctx context.Context, param_id int, amount float64, accountID int) error
	GetTransactions(ctx context.Context, account_id int) (modul.TransactionResponse, error)
}

type Repository struct {
	Accounts
	Transactions
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Accounts:     NewAccountRepistory(db),
		Transactions: NewTransactionPostgres(db),
	}
}
