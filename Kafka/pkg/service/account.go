package service

import (
	"context"
	"errors"
	"fmt"
	"go_kafka/modul.go"
	"go_kafka/pkg/repistory"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type AccountService struct {
	repo        repistory.Accounts
	redisClient *redis.Client
}

func NewAccountRepistory(repo repistory.Accounts, redisClient *redis.Client) *AccountService {
	return &AccountService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, user modul.Account) (int, error) {
	id, err := s.repo.CraeteAccount(ctx, user)
	if err != nil {
		return 0, err
	}
	key := fmt.Sprintf("accounts:%d", id)
	err = s.redisClient.HSet(ctx, key, "balance", user.Balance).Err()
	if err != nil {
		return 0, err
	}
	err = s.redisClient.Expire(ctx, key, 10*time.Hour).Err()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, account_id int) error {
	err := s.repo.DeleteAccount(ctx, account_id)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("accounts:%d", account_id)
	_, err = s.redisClient.Del(ctx, key).Result()
	if err != nil {
		return errors.New("redisten silme islemi basarisiz")
	}
	return nil
}

func (s *AccountService) GetAccountById(ctx context.Context, account_id int) (interface{}, error) {
	key := fmt.Sprintf("accounts:%d", account_id)
	balance, err := s.redisClient.HGet(ctx, key, "balance").Result()
	if err == redis.Nil {
		data, err := s.repo.GetAccountById(ctx, account_id)
		if err != nil {
			return nil, err
		}

		if data.DeletedAt != nil {
			return nil, errors.New("account is not found")
		}
		err = s.redisClient.HSet(ctx, key, "balance", data.Balance, 12*time.Hour).Err()
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"id":      data.Id,
			"balance": data.Balance,
		}, nil
	} else if err != nil {
		return nil, err
	}
	balance1, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return 0, err
	}
	return map[string]interface{}{
		"id":      account_id,
		"balance": balance1,
	}, nil
}

func (s *AccountService) GetAllAccounts(ctx context.Context) ([]map[string]interface{}, error) {
	key := "accounts"
	accountIDs, err := s.redisClient.SMembers(ctx, key).Result()
	if err == redis.Nil || len(accountIDs) == 0 {
		accounts, err := s.repo.GetAllAccounts(ctx)
		if err != nil {
			return nil, err
		}
		for _, account := range accounts {
			err := s.redisClient.HSet(ctx, fmt.Sprintf("accounts:%d", account.Id), "balance", account.Balance).Err()
			if err != nil {
				return nil, fmt.Errorf("hesap Redis'e eklenemedi: %v", err)
			}
			s.redisClient.SAdd(ctx, key, account.Id)
		}
		var liste []map[string]interface{}
		for _, account := range accounts {
			accountData := map[string]interface{}{
				"id":      account.Id,
				"balance": account.Balance,
			}
			liste = append(liste, accountData)
		}

		return liste, nil
	}
	var liste []map[string]interface{}
	for _, accountID := range accountIDs {
		balance, err := s.redisClient.HGet(ctx, fmt.Sprintf("accounts:%s", accountID), "balance").Result()
		if err != nil {
			return nil, fmt.Errorf("hesap bilgisi Redis'ten alınamadı: %v", err)
		}
		accountData := map[string]interface{}{
			"id":      accountID,
			"balance": balance,
		}
		liste = append(liste, accountData)
	}
	s.redisClient.Expire(ctx, key, 24*time.Hour)
	return liste, nil
}
