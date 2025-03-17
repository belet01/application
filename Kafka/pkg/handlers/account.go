package handler

import (
	"context"
	"go_kafka/modul.go"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Account
// @Tags api
// @Description Create a new account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body modul.Account true "Account info"
// @Success 200 {object} ResponseId
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts [post]
func (h *Handler) createAccount(c *gin.Context) {
	var input modul.Account
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	id, err := h.services.Accounts.CreateAccount(context.Background(), input)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	c.JSON(200, ResponseId{Id: id})
}

// @Summary Delete Account
// @Tags api
// @Description Delete an account by ID
// @ID delete-account-by-id
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} Message
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id} [delete]
func (h *Handler) deleteAccount(c *gin.Context) {
	ctx := context.Background()
	param := c.Param("account_id")

	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	err = h.services.Accounts.DeleteAccount(ctx, id)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	c.JSON(200, Message{Messag: "Successfully deleted"})
}

// @Summary Get Account
// @Tags api
// @Description Get account details by ID
// @ID get-account-by-id
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} modul.Account
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts/{account_id} [get]
func (h *Handler) getAccountById(c *gin.Context) {
	ctx := context.Background()
	param := c.Param("account_id")

	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	account, err := h.services.Accounts.GetAccountById(ctx, id)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	c.JSON(200, account)
}

// @Summary Get All Accounts
// @Tags api
// @Description Get a list of all accounts
// @ID get-account-all
// @Accept json
// @Produce json
// @Success 200 {array} modul.Account
// @Failure 400,404 {object} Message
// @Failure 500 {object} Message
// @Router /api/accounts [get]
func (h *Handler) getAccountAll(c *gin.Context) {
	ctx := context.Background()
	accounts, err := h.services.Accounts.GetAllAccounts(ctx)
	if err != nil {
		c.JSON(400, Message{Messag: err.Error()})
		return
	}

	c.JSON(200, accounts)
}
