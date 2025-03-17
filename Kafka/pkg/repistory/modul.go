package repistory

import (
	"context"
	"errors"
	"go_kafka/modul.go"
)

func (t *TransactionPostgres) CheckAccountActive(ctx context.Context, accountID int) (bool, error) {
	var isDeleted modul.Transactions
	query := "SELECT deleted_at FROM accounts WHERE id = $1"
	err := t.db.QueryRow(ctx, query, accountID).Scan(&isDeleted.DeletedAt)
	if err != nil {
		return false, err
	}
	if isDeleted.DeletedAt != nil {
		return false, errors.New("bu hesap silinmish")
	}
	return true, nil
}
