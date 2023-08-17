package api

import (
	"database/sql"
	"net/http"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountReq struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Currency  string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	arg := db.CreateAccountParams{
		OwnerName: req.OwnerName,
		Currency:  req.Currency,
		Balance:   0,
	}
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		// cast the pq error
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, helpers.ErrorResp(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountReq struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	account, err := server.store.GetAccountById(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, helpers.ErrorResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsReq struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req getAccountsReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
