package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferReq struct {
	FromAccountId int64  `json:"from_account_id" binding:"required"`
	ToAccountId   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR EGP CAD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	// check the currency matching
	if !server.validateAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}
	if !server.validateAccount(ctx, req.FromAccountId, req.Currency) {
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

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) bool {

	account, err := server.store.GetAccountById(ctx, accountId)

	if err != nil {
		// check if account not found
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, helpers.ErrorResp(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return false
	}

	// check the currency
	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch: %s vs %s", accountId, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return false
	}

	return true
}
