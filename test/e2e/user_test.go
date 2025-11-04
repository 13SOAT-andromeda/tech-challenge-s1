package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/session"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	var loginData handlers.LoginRequest
	var loginResponse response.BaseResponse[session.LoginOutput]
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

	loginData.Email = cfg.AdminUser.Email
	loginData.Password = cfg.AdminUser.Password

	loginResp := LoginRequest(t, loginData, apiUrl, &loginResponse)

	defer loginResp.Body.Close()

	require.NoError(t, err, "Failed to parse response body")

	assert.Equal(t, loginResp.StatusCode, http.StatusOK)
	assert.True(t, loginResponse.Success)

	t.Run("should create a valid user and delete after that", func(t *testing.T) {
		var userInput handlers.CreateUserRequest
		var userReponse response.BaseResponse[domain.User]
		var userDeleteResponse response.BaseResponse[any]

		userInput.Name = "João Marcos"
		userInput.Email = "marcosjoao1@gmail.com"
		userInput.Role = "administrator"
		userInput.Contact = "11979664877"
		userInput.Password = "Cassandra@123!"
		userInput.Address = "Rua Diamantina"
		userInput.AddressNumber = "430"
		userInput.City = "São Paulo"
		userInput.Neighborhood = "Vila Maria"
		userInput.Country = "Brasil"
		userInput.ZipCode = "02170150"

		payload, err := BuildBody(userInput)

		require.NoError(t, err, "failed on build payload request")

		resp, err := NewAuthenticatedReq("POST", apiUrl+"/users", payload, loginResponse.Data.AccessToken)

		require.NoError(t, err, "Failed on make a create user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &userReponse)

		require.NoError(t, err, "Failed to parse refresh body response")

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.True(t, userReponse.Success)
		assert.Equal(t, userInput.Name, userReponse.Data.Name)

		userId := strconv.Itoa(int(userReponse.Data.ID))

		resp, err = NewAuthenticatedReq("DELETE", apiUrl+"/users/"+userId, payload, loginResponse.Data.AccessToken)

		require.NoError(t, err, "Failed on make a delete user request")

		err = ParseBody(resp, &userDeleteResponse)

		require.NoError(t, err, "Failed to parse delete user body response")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, userDeleteResponse.Message, "User deleted")
	})
}
