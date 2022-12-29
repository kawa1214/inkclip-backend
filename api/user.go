package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/inkclip/backend/db/sqlc"
	"github.com/inkclip/backend/token"
	"github.com/inkclip/backend/util"
	"github.com/lib/pq"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// @Param request body api.registerRequest true "query params"
// @Success 200 {} {}
// @Router /register [post]
// @Tags user
func (server *Server) register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token := uuid.New()

	arg := db.CreateTemporaryUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Token:          token.String(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}
	tmpUser, err := server.store.CreateTemporaryUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// TODO: メールを送る
	print(tmpUser.Token)

	ctx.JSON(http.StatusOK, gin.H{})
}

type vertifyRequest struct {
	Email string `json:"email" form:"email" binding:"required"`
	Token string `json:"token" form:"token" binding:"required"`
}

// @Param request query api.vertifyRequest true "query params"
// @Success 200 {object} api.userResponse
// @Router /verify [get]
// @Tags user
func (server *Server) verify(ctx *gin.Context) {
	var req vertifyRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: emailとtokenで認証する
	getTmpUserArg := db.GetTemporaryUserByEmailAndTokenParams{
		Email: req.Email,
		Token: req.Token,
	}
	tmpUser, err := server.store.GetTemporaryUserByEmailAndToken(ctx, getTmpUserArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if tmpUser.ExpiresAt.Before(time.Now()) {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          tmpUser.Email,
		HashedPassword: tmpUser.HashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userResponse struct {
	ID                uuid.UUID `json:"id" binding:"required"`
	Email             string    `json:"email" binding:"required"`
	PasswordChangedAt time.Time `json:"password_changed_at" binding:"required"`
	CreatedAt         time.Time `json:"created_at" binding:"required"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// @Param request body api.createUserRequest true "query params"
// @Success 200 {object} api.userResponse
// @Router /users [post]
// @Tags user
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

// UUID support: https://github.com/gin-gonic/gin/pull/3045
type getUserRequest struct {
	ID string `uri:"id" binding:"required"`
}

// TODO: delete
func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// TODO: delete
func (server *Server) listUser(ctx *gin.Context) {
	var req listUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserRedirectResponse struct {
	SessionID             uuid.UUID    `json:"session_id" binding:"required"`
	AccessToken           string       `json:"access_token" binding:"required"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at" binding:"required"`
	RefreshToken          string       `json:"refresh_token" binding:"required"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at" binding:"required"`
	User                  userResponse `json:"user" binding:"required"`
}

// @Param request body api.loginUserRequest true "query params"
// @Success 200 {object} api.loginUserRedirectResponse
// @Router /users/login [post]
// @Tags user
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserRedirectResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// @Param request body api.renewAccessTokenRequest true "query params"
// @Success 200 {object} api.renewAccessTokenResponse
// @Router /users/renew_access [post]
// @Tags user
func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.UserID != refreshPayload.UserID {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.UserID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// @Success 200 {object} api.userResponse
// @Router /users/me [get]
// @Tags user
// @Security AccessToken
func (server *Server) getMe(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}
