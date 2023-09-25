package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	"github.com/gin-gonic/gin"
)

type renewTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type renewTokenResponse struct {
	AccessToken         string    `json:"access_token"`
	AccessTokenExpireAt time.Time `json:"access_token_expire_at"`
}

func (server *Server) requestNewAccessToken(ctx *gin.Context) {
	var req renewTokenReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	refreshPayload, err := server.tokenCreator.Verify(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.ErrorResp(err))
		return
	}

	// find the user in the database
	session, err := server.store.GetSessionById(ctx, refreshPayload.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, helpers.ErrorResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	if session.IsBlocked {
		err = fmt.Errorf("Session has been blocked!")
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
	}

	if session.Username != refreshPayload.Username {
		err = fmt.Errorf("Session doesn't belong to this user!")
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
	}

	if time.Now().After(session.ExpiredAt) {
		err = fmt.Errorf("Session has been expired before!")
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
	}

	// create the token
	token, payloadAccess, err := server.tokenCreator.Create(session.Username, server.config.TokenExpireDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	resp := renewTokenResponse{
		AccessToken:         token,
		AccessTokenExpireAt: payloadAccess.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
