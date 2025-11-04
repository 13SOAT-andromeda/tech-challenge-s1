package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
)

func SetupTest() (*config.Config, error) {

	err := initProjectRoot()

	if err != nil {
		return nil, err
	}

	cfg, err := config.Init()

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initProjectRoot() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")

	if err := os.Chdir(projectRoot); err != nil {
		return err
	}

	return nil
}

func NewUnauthenticatedReq(method, url string, body *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	return client.Do(request)
}

func NewAuthenticatedReq(method, url string, body *bytes.Buffer, token string) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)

	return client.Do(request)
}

func GetApiUrl(cfg config.Config) string {
	return cfg.Http.Url + ":" + cfg.Http.Port
}

func ParseBody[T any](body []byte, data T) {
	bodyString := string(body)
	json.Unmarshal([]byte(bodyString), &data)
}
