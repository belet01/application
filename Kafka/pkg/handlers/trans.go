package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go_kafka/modul.go"

	"github.com/gin-gonic/gin"
)

// @Summary Deposit BALANCE
// @Tags Transactions
// @Description Deposit money into an account
// @ID deposit-money
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Param input body modul.Account true "Deposit Details"
// @Success 200 {object} Message
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id}/deposit [post]
func (h *Handler) deposit(c *gin.Context) {
	var input modul.Account
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, Message{
			Messag: err.Error(),
		})
		return
	}

	ctx := context.Background()
	param := c.Param("account_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{
			Messag: err.Error(),
		})
		return
	}
	err = h.services.Transactions.DepositAccount(ctx, id, input)
	if err != nil {
		c.JSON(400, Message{
			Messag: err.Error(),
		})
		return
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- h.services.Transactions.StartDepositConsumer(ctx)
	}()
	if err := <-errCh; err != nil {
		c.JSON(400, Message{Messag: fmt.Sprintf("Hata oluştu: %v", err)})
		return
	}
	c.JSON(200, Message{
		Messag: "Successfully",
	})
}

// @Summary Withdraw Balance
// @Tags Transactions
// @Description Withdraw money from an account
// @ID withdraw-money
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Param input body modul.Account true "Withdraw Details"
// @Success 200 {object} Message
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id}/withdraw [post]

func (h *Handler) withdraw(c *gin.Context) {
	var input modul.Account
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	ctx := context.Background()
	param := c.Param("account_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	err = h.services.Transactions.WithdrawAccount(ctx, id, input)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- h.services.Transactions.StartWidthdrawConsumer(ctx)
	}()
	if err := <-errCh; err != nil {
		c.JSON(400, Message{Messag: fmt.Sprintf("Hata oluştu: %v", err)})
		return
	}
	c.JSON(200, Message{Messag: "Successfully!"})
}

// @Summary Transfer Balance
// @Tags Transactions
// @Description Transfer money between accounts
// @ID transfer-money
// @Accept json
// @Produce json
// @Param account_id path int true "Sender Account ID"
// @Param input body modul.Account true "Transfer Details"
// @Success 200 {object} Message
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id}/transfer [post]
func (h *Handler) transfer(c *gin.Context) {
	var input modul.Account
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	ctx := context.Background()
	param := c.Param("account_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	if input.Id == id {
		c.JSON(400, Message{
			Messag: "Kendi hesabınıza para göndermek için deposit işlemi yapmalısınız",
		})
		return
	}

	err = h.services.Transactions.TransferAccount(ctx, id, input.Id, input)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	// Consumer'ı Go Routine ile başlat
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh) // Kanalı kapatmayı unutma
		errCh <- h.services.Transactions.StartTransferConsumer(ctx)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			c.JSON(400, Message{Messag: fmt.Sprintf("Hata oluştu: %v", err)})
			return
		}
	case <-time.After(2000 * time.Second):
		c.JSON(400, Message{Messag: "Kafka'dan yanıt gelmedi, işlem zaman aşımına uğradı"})
		return
	}

	c.JSON(200, Message{Messag: "Successfully!"})
}

// @Summary Get Transactions
// @Tags Transactions
// @Description Get transaction history of an account
// @ID get-transaction-history
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {array} modul.Transactions
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id}/transactions [get]
func (h *Handler) getTransaction(c *gin.Context) {
	ctx := context.Background()
	param := c.Param("account_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}
	transactions, err := h.services.Transactions.GetTransaction(ctx, id)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}
	c.JSON(200, transactions)
}
