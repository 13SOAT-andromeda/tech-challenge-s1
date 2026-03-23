package e2e

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
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

	// Identity headers injected by the Lambda Authorizer for the admin user.
	adminEmail := cfg.AdminUser.Email
	adminRole := "administrator"
	adminID := "1"

	var createdUserID uint

	t.Run("should create a valid user", func(t *testing.T) {
		var userInput handlers.CreateUserRequest
		var userResponse domain.User

		userInput.Name = "João Marcos"
		userInput.Email = "marcosjoao" + strconv.FormatInt(time.Now().Unix(), 10) + "@gmail.com"
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

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a create user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &userResponse)

		require.NoError(t, err, "Failed to parse create user response body")

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, userInput.Name, userResponse.Name)
		assert.Equal(t, userInput.Email, userResponse.Email)
		assert.Equal(t, userInput.Role, userResponse.Role)
		assert.Equal(t, userInput.Contact, userResponse.Contact)
		assert.NotZero(t, userResponse.ID)

		createdUserID = userResponse.ID
	})

	t.Run("should fail to create user with invalid request", func(t *testing.T) {
		var userInput handlers.CreateUserRequest
		var errorResponse map[string]interface{}

		userInput.Name = ""
		userInput.Email = "invalid-email"
		userInput.Password = "weak"

		payload, err := BuildBody(userInput)

		require.NoError(t, err, "failed on build payload request")

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a create user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to create user with duplicate email", func(t *testing.T) {
		var userInput handlers.CreateUserRequest
		var errorResponse map[string]interface{}

		userInput.Name = "João Duplicado"
		userInput.Email = cfg.AdminUser.Email
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

		resp, err := NewIdentifiedReq("POST", apiUrl+"/users", payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a create user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should get user by ID successfully", func(t *testing.T) {
		var userResponse domain.User

		userId := strconv.Itoa(int(createdUserID))

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a get user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &userResponse)

		require.NoError(t, err, "Failed to parse get user response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, createdUserID, userResponse.ID)
		assert.NotEmpty(t, userResponse.Name)
		assert.NotEmpty(t, userResponse.Email)
	})

	t.Run("should fail to get user with invalid ID", func(t *testing.T) {
		var errorResponse map[string]interface{}

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/invalid", nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a get user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to get non-existent user", func(t *testing.T) {
		var errorResponse map[string]interface{}

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/999999", nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a get user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should search users successfully", func(t *testing.T) {
		var usersResponse []domain.User

		params := url.Values{}
		params.Add("name", "João")
		searchUrl := apiUrl + "/users?" + params.Encode()

		resp, err := NewIdentifiedReq("GET", searchUrl, nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a search users request")

		defer resp.Body.Close()

		err = ParseBody(resp, &usersResponse)

		require.NoError(t, err, "Failed to parse search users response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, usersResponse)
		assert.GreaterOrEqual(t, len(usersResponse), 1)
	})

	t.Run("should search users without parameters", func(t *testing.T) {
		var usersResponse []domain.User

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users", nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a search users request")

		defer resp.Body.Close()

		err = ParseBody(resp, &usersResponse)

		require.NoError(t, err, "Failed to parse search users response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, usersResponse)
	})

	t.Run("should update user successfully", func(t *testing.T) {
		var updateInput handlers.UpdateUserRequest
		var updateResponse map[string]interface{}

		updateInput.Name = "João Marcos Atualizado"
		updateInput.Contact = "11988888888"
		updateInput.Address = "Rua Atualizada"
		updateInput.AddressNumber = "500"
		updateInput.City = "Rio de Janeiro"
		updateInput.Neighborhood = "Centro"
		updateInput.Country = "Brasil"
		updateInput.ZipCode = "20000000"

		payload, err := BuildBody(updateInput)

		require.NoError(t, err, "failed on build payload request")

		userId := strconv.Itoa(int(createdUserID))

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/"+userId, payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make an update user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &updateResponse)

		require.NoError(t, err, "Failed to parse update user response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, updateResponse["message"], "updated")

		// Verify the update
		var userResponse domain.User
		resp, err = NewIdentifiedReq("GET", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)
		require.NoError(t, err)
		defer resp.Body.Close()

		err = ParseBody(resp, &userResponse)
		require.NoError(t, err)

		assert.Equal(t, updateInput.Name, userResponse.Name)
		assert.Equal(t, updateInput.Contact, userResponse.Contact)
	})

	t.Run("should fail to update user with invalid ID", func(t *testing.T) {
		var updateInput handlers.UpdateUserRequest
		var errorResponse map[string]interface{}

		updateInput.Name = "Test"

		payload, err := BuildBody(updateInput)

		require.NoError(t, err, "failed on build payload request")

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/invalid", payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make an update user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail to update non-existent user", func(t *testing.T) {
		var updateInput handlers.UpdateUserRequest
		var errorResponse map[string]interface{}

		updateInput.Name = "Test User"

		payload, err := BuildBody(updateInput)

		require.NoError(t, err, "failed on build payload request")

		resp, err := NewIdentifiedReq("PUT", apiUrl+"/users/999999", payload, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make an update user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should delete user successfully", func(t *testing.T) {
		var deleteResponse map[string]interface{}

		userId := strconv.Itoa(int(createdUserID))

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a delete user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &deleteResponse)

		require.NoError(t, err, "Failed to parse delete user response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, deleteResponse["message"], "deleted")
	})

	t.Run("should fail to delete user with invalid ID", func(t *testing.T) {
		var errorResponse map[string]interface{}

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/invalid", nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a delete user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should delete non-existent user without error", func(t *testing.T) {
		var deleteResponse map[string]interface{}

		resp, err := NewIdentifiedReq("DELETE", apiUrl+"/users/999999", nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a delete user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &deleteResponse)

		require.NoError(t, err, "Failed to parse delete response body")

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, deleteResponse["message"], "deleted")
	})

	t.Run("should verify deleted user is not found", func(t *testing.T) {
		var errorResponse map[string]interface{}

		userId := strconv.Itoa(int(createdUserID))

		resp, err := NewIdentifiedReq("GET", apiUrl+"/users/"+userId, nil, adminID, adminEmail, adminRole)

		require.NoError(t, err, "Failed on make a get user request")

		defer resp.Body.Close()

		err = ParseBody(resp, &errorResponse)

		require.NoError(t, err, "Failed to parse error response body")

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
