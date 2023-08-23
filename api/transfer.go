package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/token"
	"github.com/gin-gonic/gin"
)

type createTransferReq struct {
	FromAccountId int64  `json:"from_account_id" binding:"required"`
	ToAccountId   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	// check the currency matching
	_, valid := server.validateAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
		return
	}
	fromAccount, valid := server.validateAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}
	// logged-in user can only transfer money from his account to others
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if payload.Username != fromAccount.OwnerName {
		ctx.JSON(http.StatusUnauthorized,
			helpers.ErrorResp(errors.New("from_account doesn't belong to the logged-in user!")),
		)
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTransaction(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {

	account, err := server.store.GetAccountById(ctx, accountId)

	if err != nil {
		// check if account not found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, helpers.ErrorResp(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return account, false
	}

	// check the currency
	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch: %s vs %s", accountId, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return account, false
	}

	return account, true
}
