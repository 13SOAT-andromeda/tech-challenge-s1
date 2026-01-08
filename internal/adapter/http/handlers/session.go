package handlers

import (
	"net/http"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	useCase ports.SessionUseCase
}

func NewSessionHandler(
	useCase ports.SessionUseCase,
) *SessionHandler {
	return &SessionHandler{
		useCase: useCase,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *SessionHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	input := ports.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.useCase.Login(ctx, input)
	if err != nil {
		response.RespondError(ctx, mapErrorToStatus(err), err.Error())
		return
	}

	response.RespondSuccess(ctx, output, "")
}

func (h *SessionHandler) Validate(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.RespondError(ctx, http.StatusUnauthorized, "Authorization header required")
		return
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		response.RespondError(ctx, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}

	token := authHeader[7:]
	input := ports.ValidateInput{
		Token: token,
	}

	output, err := h.useCase.Validate(ctx, input)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, output, "")
}

func (h *SessionHandler) Refresh(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	input := ports.RefreshInput{
		RefreshToken: req.RefreshToken,
	}

	output, err := h.useCase.Refresh(ctx, input)
	if err != nil {
		response.RespondError(ctx, mapErrorToStatus(err), err.Error())
		return
	}

	response.RespondSuccess(ctx, output, "")
}

func (h *SessionHandler) Logout(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	input := ports.LogoutInput{
		RefreshToken: req.RefreshToken,
	}

	err := h.useCase.Logout(ctx, input)
	if err != nil {
		response.RespondError(ctx, mapErrorToStatus(err), err.Error())
		return
	}

	response.RespondSuccess(ctx, "Logged out successfully", "")
}

func mapErrorToStatus(err error) int {
	switch err {
	case services.ErrUserNotFound:
		return http.StatusUnauthorized
	case services.ErrSessionInvalid:
		return http.StatusUnauthorized
	case services.ErrSessionRefreshTokenEmpty:
		return http.StatusBadRequest
	case services.ErrSessionExpiresAtPast:
		return http.StatusBadRequest
	case services.ErrSessionUserIDInvalid:
		return http.StatusBadRequest
	case services.ErrSessionIDInvalid:
		return http.StatusBadRequest
	case services.ErrSessionNil:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
