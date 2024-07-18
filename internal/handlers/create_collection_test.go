package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/skillBuilder/internal/mocks"
)

type result struct {
	ID string `json:"id"`
}

func TestCreateCollection(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	mockStorage := &mocks.MockStorage{
		CreateCollectionFunc: func(userID string, name string, description string) (string, error) {
			return "test_collection_id", nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(Collection{Name: "Test Collection", Description: "Test Description"})
	req, err := http.NewRequest("POST", "/api/collection/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.CreateCollection(r, req)

	rr := httptest.NewRecorder()

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestCreateCollection_CreateError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	mockStorage := &mocks.MockStorage{
		CreateCollectionFunc: func(userID string, name string, description string) (string, error) {
			return "", errors.New("create collection error")
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(Collection{Name: "Test Collection", Description: "Test Description"})
	req, err := http.NewRequest("POST", "/api/collection/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.CreateCollection(r, req)

	assert.Equal(t, http.StatusUnauthorized, r.Code)
}
