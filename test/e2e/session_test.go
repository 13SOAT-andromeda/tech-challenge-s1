package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/session"
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
		var baseResponse response.BaseResponse[session.LoginOutput]

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		bodyBytes, resp := loginRequest(t, loginData, apiUrl)

		err = ParseBody(bodyBytes, &baseResponse)

		require.NoError(t, err, "Failed to parse response body")

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
		assert.True(t, baseResponse.Data.User.Active)
	})

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[any]

		loginData.Email = "invalidmail@test.com"
		loginData.Password = "invalidpass"

		bodyBytes, resp := loginRequest(t, loginData, apiUrl)

		err = ParseBody(bodyBytes, &baseResponse)

		require.NoError(t, err, "Failed to parse response body")

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
		assert.False(t, baseResponse.Success)
		assert.Contains(t, baseResponse.Message, "usuário não encontrado")
	})

	t.Run("should fail with invalid request", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[any]

		loginData.Email = "invalidmail@test.com"

		bodyBytes, resp := loginRequest(t, loginData, apiUrl)

		err = ParseBody(bodyBytes, &baseResponse)

		require.NoError(t, err, "Failed to parse response body")

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.False(t, baseResponse.Success)
	})

	t.Run("should validate token results in successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[session.LoginOutput]
		var baseResponseValidate response.BaseResponse[session.ValidateOutput]

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		bodyBytes, resp := loginRequest(t, loginData, apiUrl)

		err = ParseBody(bodyBytes, &baseResponse)

		require.NoError(t, err, "Failed to parse response body")

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
		assert.True(t, baseResponse.Data.User.Active)

		resp, err := NewAuthenticatedReq("GET", apiUrl+"/sessions/validate", nil, baseResponse.Data.AccessToken)

		require.NoError(t, err, "Failed on make a validate request")

		bodyBytes, err = io.ReadAll(resp.Body)

		require.NoError(t, err, "Failed to parse validate body response")

		defer resp.Body.Close()

		err = ParseBody(bodyBytes, &baseResponseValidate)

		require.NoError(t, err, "Failed to parse validate response body")

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.True(t, baseResponseValidate.Success)
		assert.True(t, baseResponseValidate.Data.Valid)
		assert.NotEmpty(t, baseResponseValidate.Data.User)
		assert.NotNil(t, baseResponseValidate.Data.User)
		assert.NotZero(t, baseResponseValidate.Data.User.ID)
		assert.Equal(t, cfg.AdminUser.Email, baseResponseValidate.Data.User.Email)
		assert.NotEmpty(t, baseResponseValidate.Data.User.Role)
		assert.True(t, baseResponseValidate.Data.User.Active)
	})
}

func loginRequest(t *testing.T, loginData handlers.LoginRequest, apiUrl string) ([]byte, *http.Response) {
	payloadJson, err := json.Marshal(loginData)

	require.NoError(t, err, "Failed to marshal login request")

	payload := bytes.NewBuffer(payloadJson)
	resp, err := NewUnauthenticatedReq("POST", apiUrl+"/sessions", payload)

	require.NoError(t, err, "Failed on make a request")

	bodyBytes, err := io.ReadAll(resp.Body)

	require.NoError(t, err, "Failed to parse body response")

	defer resp.Body.Close()

	return bodyBytes, resp
}
