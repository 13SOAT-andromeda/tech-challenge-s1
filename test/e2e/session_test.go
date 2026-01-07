package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	cfg, err := SetupTest()

	if err != nil {
		t.Skipf("Error on setup config variables: %v", err)
	}

	apiUrl := GetApiUrl(*cfg)

	healthResp, err := http.Get(apiUrl + "/health")
	if err != nil {
		t.Skipf("Application not running at %s, skipping E2E tests: %v", apiUrl, err)
		return
	}
	healthResp.Body.Close()

	t.Run("should login successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[ports.LoginOutput]

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponse.Success)
		assert.NotEmpty(t, baseResponse.Data.AccessToken)
		assert.NotEmpty(t, baseResponse.Data.RefreshToken)
		assert.NotEmpty(t, baseResponse.Data.ExpiresIn)
		assert.NotEmpty(t, baseResponse.Data.User)
		assert.NotNil(t, baseResponse.Data.User)
		assert.NotZero(t, baseResponse.Data.User.ID)
		assert.Equal(t, cfg.AdminUser.Email, baseResponse.Data.User.Email)
		assert.NotEmpty(t, baseResponse.Data.User.Role)
	})

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[any]

		loginData.Email = "invalidmail@test.com"
		loginData.Password = "invalidpass"

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.False(t, baseResponse.Success)
		assert.Contains(t, baseResponse.Message, "usuário não encontrado")
	})

	t.Run("should fail with invalid request", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[any]

		loginData.Email = "invalidmail@test.com"

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.False(t, baseResponse.Success)
	})

	t.Run("should validate token results in successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[ports.LoginOutput]
		var baseResponseValidate response.BaseResponse[ports.ValidateOutput]

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponse.Success)
		assert.NotEmpty(t, baseResponse.Data.AccessToken)
		assert.NotEmpty(t, baseResponse.Data.RefreshToken)
		assert.NotEmpty(t, baseResponse.Data.ExpiresIn)
		assert.NotEmpty(t, baseResponse.Data.User)
		assert.NotNil(t, baseResponse.Data.User)
		assert.NotZero(t, baseResponse.Data.User.ID)
		assert.Equal(t, cfg.AdminUser.Email, baseResponse.Data.User.Email)
		assert.NotEmpty(t, baseResponse.Data.User.Role)

		resp, err := NewAuthenticatedReq("GET", apiUrl+"/sessions/validate", nil, baseResponse.Data.AccessToken)

		require.NoError(t, err, "Failed on make a validate request")

		err = ParseBody(resp, &baseResponseValidate)

		require.NoError(t, err, "Failed to parse validate response body")

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponseValidate.Success)
		assert.True(t, baseResponseValidate.Data.Valid)
		assert.NotEmpty(t, baseResponseValidate.Data.User)
		assert.NotNil(t, baseResponseValidate.Data.User)
		assert.NotZero(t, baseResponseValidate.Data.User.ID)
		assert.Equal(t, cfg.AdminUser.Email, baseResponseValidate.Data.User.Email)
		assert.NotEmpty(t, baseResponseValidate.Data.User.Role)
	})

	t.Run("should invalid token return an error", func(t *testing.T) {
		var baseResponseValidate response.BaseResponse[any]

		resp, err := NewAuthenticatedReq("GET", apiUrl+"/sessions/validate", nil, "invalidtoken")

		require.NoError(t, err, "Failed on make a validate request")

		defer resp.Body.Close()

		err = ParseBody(resp, &baseResponseValidate)

		require.NoError(t, err, "Failed to parse validate response body")

		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
		assert.False(t, baseResponseValidate.Success)
		assert.Contains(t, baseResponseValidate.Message, "token is malformed")
	})

	t.Run("should refresh token route results in successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var refreshData handlers.RefreshRequest
		var baseResponse response.BaseResponse[ports.LoginOutput]
		var baseResponseRefresh response.BaseResponse[ports.RefreshOutput]

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponse.Success)
		assert.NotEmpty(t, baseResponse.Data.AccessToken)
		assert.NotEmpty(t, baseResponse.Data.RefreshToken)
		assert.NotEmpty(t, baseResponse.Data.ExpiresIn)
		assert.NotEmpty(t, baseResponse.Data.User)
		assert.NotNil(t, baseResponse.Data.User)
		assert.NotZero(t, baseResponse.Data.User.ID)
		assert.Equal(t, cfg.AdminUser.Email, baseResponse.Data.User.Email)
		assert.NotEmpty(t, baseResponse.Data.User.Role)

		refreshData.RefreshToken = baseResponse.Data.RefreshToken

		payload, err := BuildBody(refreshData)

		require.NoError(t, err, "Failed to build refresh request")

		resp, err = NewUnauthenticatedReq("POST", apiUrl+"/sessions/refresh", payload)

		require.NoError(t, err, "Failed on make a refresh request")

		err = ParseBody(resp, &baseResponseRefresh)

		require.NoError(t, err, "Failed to parse refresh response body")

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponseRefresh.Success)
		assert.NotEmpty(t, baseResponseRefresh.Data.AccessToken)
		assert.NotEmpty(t, baseResponseRefresh.Data.RefreshToken)
		assert.NotZero(t, baseResponseRefresh.Data.ExpiresIn)
	})

	t.Run("should invalid refresh token return an error", func(t *testing.T) {
		var baseResponseRefresh response.BaseResponse[any]
		var refreshData handlers.RefreshRequest

		refreshData.RefreshToken = "invalidtoken"

		body, err := BuildBody(refreshData)

		require.NoError(t, err, "Failed to build refresh validate body")

		resp, err := NewUnauthenticatedReq("POST", apiUrl+"/sessions/refresh", body)

		require.NoError(t, err, "Failed on make a refresh validate request")

		defer resp.Body.Close()

		err = ParseBody(resp, &baseResponseRefresh)

		require.NoError(t, err, "Failed to parse refresh validate response body")

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.False(t, baseResponseRefresh.Success)
		assert.Contains(t, baseResponseRefresh.Message, "sessão inválida ou expirada")
	})

	t.Run("should logout successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[ports.LoginOutput]
		var baseResponseLogout response.BaseResponse[any]
		var logoutData handlers.RefreshRequest

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		resp := LoginRequest(t, loginData, apiUrl, &baseResponse)

		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponse.Success)
		assert.NotEmpty(t, baseResponse.Data.RefreshToken)

		logoutData.RefreshToken = baseResponse.Data.RefreshToken

		payload, err := BuildBody(logoutData)

		require.NoError(t, err, "Failed to build logout request")

		resp, err = NewUnauthenticatedReq("DELETE", apiUrl+"/sessions/logout", payload)

		require.NoError(t, err, "Failed on make a logout request")

		err = ParseBody(resp, &baseResponseLogout)

		require.NoError(t, err, "Failed to parse logout response body")

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponseLogout.Success)
		assert.Equal(t, "Logged out successfully", baseResponseLogout.Data)
	})

	t.Run("should invalid refresh token return an error on logout", func(t *testing.T) {
		var baseResponseLogout response.BaseResponse[any]
		var logoutData handlers.RefreshRequest

		logoutData.RefreshToken = "invalidtoken"

		payload, err := BuildBody(logoutData)

		require.NoError(t, err, "Failed to build logout request")

		resp, err := NewUnauthenticatedReq("DELETE", apiUrl+"/sessions/logout", payload)

		require.NoError(t, err, "Failed on make a logout request")

		defer resp.Body.Close()

		err = ParseBody(resp, &baseResponseLogout)

		require.NoError(t, err, "Failed to parse logout response body")

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.False(t, baseResponseLogout.Success)
		assert.Contains(t, baseResponseLogout.Message, "sessão inválida ou expirada")
	})
}

func LoginRequest[T any](t *testing.T, loginData handlers.LoginRequest, apiUrl string, output *response.BaseResponse[T]) *http.Response {
	payload, err := BuildBody(loginData)

	require.NoError(t, err, "Failed to build login request")

	resp, err := NewUnauthenticatedReq("POST", apiUrl+"/sessions", payload)

	require.NoError(t, err, "Failed on make a request")

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	require.NoError(t, err, "Failed to parse body response")

	bodyString := string(bodyBytes)

	err = json.Unmarshal([]byte(bodyString), &output)

	require.NoError(t, err, "Failed to unmarshal body response")

	return resp
}
