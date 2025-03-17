package service

import (
	"context"
	"go_kafka/modul.go"
	"go_kafka/pkg/repistory"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-redis/redis/v8"
)

type Accounts interface {
	CreateAccount(ctx context.Context, user modul.Account) (int, error)
	DeleteAccount(ctx context.Context, account_id int) error
	GetAccountById(ctx context.Context, account_id int) (interface{}, error)
	GetAllAccounts(ctx context.Context) ([]map[string]interface{}, error)
}

type Transactions interface {
	DepositAccount(ctx context.Context, accountID int, USER modul.Account) error
	StartDepositConsumer(ctx context.Context) error
	WithdrawAccount(ctx context.Context, accountID int, user modul.Account) error
	StartWidthdrawConsumer(ctx context.Context) error
	TransferAccount(ctx context.Context, param_id, accountID int, user modul.Account) error
	StartTransferConsumer(ctx context.Context) error
	GetTransaction(ctx context.Context, account_id int) (modul.TransactionResponse, error)
}

type Service struct {
	Accounts
	Transactions
}

func NewService(repo repistory.Repository, redisClient *redis.Client, prod *kafka.Producer, cons *kafka.Consumer) *Service {
	return &Service{
		Accounts:     NewAccountRepistory(repo.Accounts, redisClient),
		Transactions: NewTransactionService(repo.Transactions, redisClient, prod, cons),
	}
}
