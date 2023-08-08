package api

import (
	"net/http"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountReq struct {
	OwnerName string `json:"owner_name" binding:"required"`
	Currency  string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq

	if err := ctx.BindJSON(&req); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
