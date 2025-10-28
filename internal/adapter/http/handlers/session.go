package handlers

import (
	"net/http"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/session"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	loginUseCase    session.LoginUseCase
	validateUseCase session.ValidateUseCase
	refreshUseCase  session.RefreshUseCase
	logoutUseCase   session.LogoutUseCase
}

func NewSessionHandler(
	loginUC session.LoginUseCase,
	validateUC session.ValidateUseCase,
	refreshUC session.RefreshUseCase,
	logoutUC session.LogoutUseCase,
) *SessionHandler {
	return &SessionHandler{
		loginUseCase:    loginUC,
		validateUseCase: validateUC,
		refreshUseCase:  refreshUC,
		logoutUseCase:   logoutUC,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login handles POST /sessions
func (h *SessionHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := session.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.loginUseCase.Execute(ctx, input)
	if err != nil {
		ctx.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Validate handles GET /sessions/validate
func (h *SessionHandler) Validate(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Extract token from "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	token := authHeader[7:]
	input := session.ValidateInput{
		Token: token,
	}

	output, err := h.validateUseCase.Execute(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Refresh handles POST /sessions/refresh
func (h *SessionHandler) Refresh(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := session.RefreshInput{
		RefreshToken: req.RefreshToken,
	}

	output, err := h.refreshUseCase.Execute(ctx, input)
	if err != nil {
		ctx.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, output)
}

// Logout handles DELETE /sessions/logout
func (h *SessionHandler) Logout(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := session.LogoutInput{
		RefreshToken: req.RefreshToken,
	}

	err := h.logoutUseCase.Execute(ctx, input)
	if err != nil {
		ctx.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// mapErrorToStatus maps domain errors to HTTP status codes
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
