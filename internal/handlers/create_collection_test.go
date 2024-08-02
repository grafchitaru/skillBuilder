package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/grafchitaru/skillBuilder/internal/models"
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
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		CreateCollectionFunc: func(userID string, name string, description string) (string, error) {
			return "test_collection_id", nil
		},
		AddCollectionToUserFunc: func(userID string, name string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.Collection{Name: "Test Collection", Description: "Test Description"})
	req, err := http.NewRequest("POST", "/api/collection/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Add the authentication cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWYwMmQwMzYtYjQ1Ny00M2ExLThmYzktNWM2NDBjM2Y3ZDJhIn0.F2NM790xbzXL6b-gpxg3xUp1G76ZHS43Gy0dZwGlmJg",
		Path:  "/",
	})
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.CreateCollection(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)
}

func TestCreateCollection_CreateError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		CreateCollectionFunc: func(userID string, name string, description string) (string, error) {
			return "", errors.New("create collection error")
		},
		AddCollectionToUserFunc: func(userID string, name string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.Collection{Name: "Test Collection", Description: "Test Description"})
	req, err := http.NewRequest("POST", "/api/collection/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Add the authentication cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWYwMmQwMzYtYjQ1Ny00M2ExLThmYzktNWM2NDBjM2Y3ZDJhIn0.F2NM790xbzXL6b-gpxg3xUp1G76ZHS43Gy0dZwGlmJg",
		Path:  "/",
	})
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.CreateCollection(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
}
