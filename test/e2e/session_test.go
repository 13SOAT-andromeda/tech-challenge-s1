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
)

func TestSession(t *testing.T) {
	cfg, err := SetupTest()
	apiUrl := GetApiUrl(*cfg)

	t.Run("should login successful", func(t *testing.T) {
		var loginData handlers.LoginRequest
		var baseResponse response.BaseResponse[session.LoginOutput]
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}

		loginData.Email = cfg.AdminUser.Email
		loginData.Password = cfg.AdminUser.Password

		payloadJson, err := json.Marshal(loginData)

		if err != nil {
			t.Fatalf("Failed on parse body to JSON: %v", err)
		}

		payload := bytes.NewBuffer(payloadJson)
		response, err := NewUnauthenticatedReq("POST", apiUrl+"/sessions", payload)

		if err != nil {
			t.Fatalf("Failed on make a request: %v", err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Failed to parse body response: %v", err)
		}
		defer response.Body.Close()

		ParseBody(bodyBytes, &baseResponse)

		assert.Equal(t, response.StatusCode, http.StatusOK)
		assert.True(t, baseResponse.Success)
		assert.NotEmpty(t, baseResponse.Data.AccessToken)
		assert.NotEmpty(t, baseResponse.Data.RefreshToken)
		assert.NotEmpty(t, baseResponse.Data.ExpiresIn)
		assert.NotEmpty(t, baseResponse.Data.User)
	})
}
