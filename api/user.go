package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AYehia0/go-bk-mst/api/helpers"
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

type userResp struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

func newUserResp(user db.User) userResp {
	return userResp{
		FullName: user.FullName,
		Email:    user.Email,
		Username: user.Username,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req userReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	// hashing the password
	hashedPassword, err := utils.GenerateHash(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
	}
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		// cast the pq error
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, helpers.ErrorResp(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	// without the password field
	userTemp := newUserResp(user)
	ctx.JSON(http.StatusOK, userTemp)
}

type userLoginReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
type loginResponse struct {
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	User                 userResp  `json:"user"`
	AccessTokenExpireAt  time.Time `json:"access_token_expire_at"`
	RefreshTokenExpireAt time.Time `json:"refresh_token_expire_at"`
	SessionId            uuid.UUID `json:"session_id"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req userLoginReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.ErrorResp(err))
		return
	}

	// find the user in the database
	user, err := server.store.GetUserByUsername(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, helpers.ErrorResp(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	// checking the password
	err = utils.ComparePasswords(req.Password, user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.ErrorResp(err))
		return
	}

	// create the token
	token, payloadAccess, err := server.tokenCreator.Create(user.Username, server.config.TokenExpireDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	// create the refresh token.
	// the refresh token should be linked to a user.
	refreshToken, payloadRefresh, err := server.tokenCreator.Create(req.Username, server.config.TokenRefreshExpireDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	// the refresh token's UUID --> ID should be returned from the payload itself but the tokenCreator.Create returns the encrypted payload hence we should also return the payload.
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           payloadRefresh.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		IpAddr:       ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpiredAt:    payloadRefresh.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.ErrorResp(err))
		return
	}

	resp := loginResponse{
		AccessToken:          token,
		AccessTokenExpireAt:  payloadAccess.ExpiredAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: payloadRefresh.ExpiredAt,
		SessionId:            session.ID,
		User:                 newUserResp(user),
	}

	ctx.JSON(http.StatusOK, resp)
}
