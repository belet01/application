package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_kafka/modul.go"
	"go_kafka/pkg/repistory"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-redis/redis/v8"
)

type TransactionServise struct {
	repo          repistory.Transactions
	redisClient   *redis.Client
	kafkaProduser *kafka.Producer
	kafkaConsumer *kafka.Consumer
}

func NewTransactionService(repo repistory.Transactions, redisClient *redis.Client, prod *kafka.Producer, cons *kafka.Consumer) *TransactionServise {
	return &TransactionServise{
		repo:          repo,
		redisClient:   redisClient,
		kafkaProduser: prod,
		kafkaConsumer: cons,
	}
}

func (t *TransactionServise) DepositAccount(ctx context.Context, accountID int, user modul.Account) error {
	del, err := t.repo.CheckAccountActive(ctx, accountID)
	if err != nil {
		return err
	}
	if !del {
		return errors.New("account is deleted")
	}
	transaction := modul.TransactionREsponse{
		AccountId:       accountID,
		Amount:          user.Balance,
		TransactionType: "deposiy",
	}
	message, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	err = t.kafkaProduser.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &[]string{"account-deposit"}[0],
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(fmt.Sprintf("%d", accountID)),
		Value: message,
	}, nil)
	if err != nil {
		return err
	}
	t.kafkaProduser.Flush(15 * 1000)
	key := fmt.Sprintf("account:%d", accountID)
	_, err = t.redisClient.Del(ctx, key).Result()
	if err != nil {
		return errors.New("redisten silme islemi basarisiz")
	}
	return nil
}

func (t *TransactionServise) StartDepositConsumer(ctx context.Context) error {
	err := t.kafkaConsumer.SubscribeTopics([]string{"account-deposit"}, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := t.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			return nil
		}
		var transaction modul.TransactionREsponse
		err = json.Unmarshal(msg.Value, &transaction)
		if err != nil {
			return err
		}
		if err = t.repo.Deposit(ctx, transaction.Amount, transaction.AccountId); err != nil {
			return err
		}

	}

}

func (t *TransactionServise) WithdrawAccount(ctx context.Context, accountID int, user modul.Account) error {
	del, err := t.repo.CheckAccountActive(ctx, accountID)
	if err != nil {
		return err
	}
	if !del {
		return err
	}

	transaction := modul.TransactionREsponse{
		AccountId:       accountID,
		Amount:          user.Balance,
		TransactionType: "withdraw",
	}
	message, err := json.Marshal(&transaction)
	if err != nil {
		return err
	}

	err = t.kafkaProduser.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &[]string{"account-withdrow"}[0],
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(fmt.Sprintf("%d", accountID)),
		Value: message,
	}, nil)
	if err != nil {
		return err
	}
	t.kafkaProduser.Flush(15 * 1000)
	key := fmt.Sprintf("account:%d", accountID)
	_, err = t.redisClient.Del(ctx, key).Result()
	if err != nil {
		return errors.New("redisten silme islemi basarisiz")
	}
	return nil
}

func (t *TransactionServise) StartWidthdrawConsumer(ctx context.Context) error {
	err := t.kafkaConsumer.SubscribeTopics([]string{"account-withdrow"}, nil)
	if err != nil {
		return fmt.Errorf("kafka'ya abone olunamadı: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := t.kafkaConsumer.ReadMessage(-1)
			if err != nil {
				return err
			}
			var transaction modul.TransactionREsponse
			if err := json.Unmarshal(msg.Value, &transaction); err != nil {
				return err
			}

			if err := t.repo.Withdraw(ctx, transaction.Amount, transaction.AccountId); err != nil {

				return err
			}

		}
	}
}

func (t *TransactionServise) TransferAccount(ctx context.Context, param_id, accountID int, user modul.Account) error {
	del, err := t.repo.CheckAccountActive(ctx, accountID)
	if err != nil {
		return err
	}
	if !del {
		return errors.New("bu hesaba transfer islemi yapamazsiniz.. Hesap silinmish")
	}
	del, err = t.repo.CheckAccountActive(ctx, param_id)
	if err != nil {
		return err
	}
	if !del {
		return errors.New("boyle bir hesap yok")
	}
	transaction := modul.TransactionREsponse{
		YourId:          param_id,
		AccountId:       accountID,
		Amount:          user.Balance,
		TransactionType: "transfer",
	}
	message, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	err = t.kafkaProduser.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &[]string{"account-transfer"}[0],
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(fmt.Sprintf("%d", accountID)),
		Value: message,
	}, nil)
	if err != nil {
		return fmt.Errorf("kafka'ya mesaj gönderilemedi: %v", err)
	}
	t.kafkaProduser.Flush(15 * 1000)

	key := fmt.Sprintf("account:%d", accountID)
	_, err = t.redisClient.Del(ctx, key).Result()
	if err != nil {
		return errors.New("redisten silme islemi basarisiz")
	}
	key1 := fmt.Sprintf("account:%d", param_id)
	_, err = t.redisClient.Del(ctx, key1).Result()
	if err != nil {
		return errors.New("redisten silme islemi basarisiz")
	}
	return nil
}

func (t *TransactionServise) StartTransferConsumer(ctx context.Context) error {
	err := t.kafkaConsumer.SubscribeTopics([]string{"account-transfer"}, nil)
	if err != nil {
		return fmt.Errorf("kafka'ya abone olunamadı: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := t.kafkaConsumer.ReadMessage(-1)
			if err != nil {
				return err
			}
			var transaction modul.TransactionREsponse
			if err := json.Unmarshal(msg.Value, &transaction); err != nil {
				return err
			}

			if err := t.repo.Transfer(ctx, transaction.YourId, transaction.Amount, transaction.AccountId); err != nil {

				return err
			}

		}
	}

}

func (t *TransactionServise) GetTransaction(ctx context.Context, account_id int) (modul.TransactionResponse, error) {
	del, err := t.repo.CheckAccountActive(ctx, account_id)
	if err != nil {
		return modul.TransactionResponse{}, err
	}
	if !del {
		return modul.TransactionResponse{}, errors.New("bu hesap silinmish")
	}
	return t.repo.GetTransactions(ctx, account_id)
}
