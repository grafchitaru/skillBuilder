package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteCollectionFromUser(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		DeleteCollectionFromUserFunc: func(userID string, collectionID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}

	r := chi.NewRouter()
	r.Delete("/api/collection/{id}/user", hc.DeleteCollectionFromUser)

	req, err := http.NewRequest("DELETE", "/api/collection/a97ae726-3859-4a1f-85ec-22452c243ac5/user", nil)
	req.Header.Set("Content-Type", "application/json")
	// Add the authentication cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWYwMmQwMzYtYjQ1Ny00M2ExLThmYzktNWM2NDBjM2Y3ZDJhIn0.F2NM790xbzXL6b-gpxg3xUp1G76ZHS43Gy0dZwGlmJg",
		Path:  "/",
	})
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteCollectionFromUser_IDNotFound(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		DeleteCollectionFromUserFunc: func(userID string, collectionID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}

	r := chi.NewRouter()
	r.Delete("/api/collection/{id}/user", hc.DeleteCollectionFromUser)

	req, err := http.NewRequest("DELETE", "/api/collection//user", nil)
	req.Header.Set("Content-Type", "application/json")
	// Add the authentication cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWYwMmQwMzYtYjQ1Ny00M2ExLThmYzktNWM2NDBjM2Y3ZDJhIn0.F2NM790xbzXL6b-gpxg3xUp1G76ZHS43Gy0dZwGlmJg",
		Path:  "/",
	})
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteCollectionFromUser_Unauthorized(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{
		DeleteCollectionFromUserFunc: func(userID string, collectionID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return "", errors.New("unauthorized")
	}

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}

	r := chi.NewRouter()
	r.Delete("/api/collection/{id}/user", hc.DeleteCollectionFromUser)

	req, err := http.NewRequest("DELETE", "/api/collection/a97ae726-3859-4a1f-85ec-22452c243ac5/user", nil)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteCollectionFromUser_DeleteError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		DeleteCollectionFromUserFunc: func(userID string, collectionID string) error {
			return errors.New("delete collection from user error")
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}

	r := chi.NewRouter()
	r.Delete("/api/collection/{id}/user", hc.DeleteCollectionFromUser)

	req, err := http.NewRequest("DELETE", "/api/collection/a97ae726-3859-4a1f-85ec-22452c243ac5/user", nil)
	req.Header.Set("Content-Type", "application/json")
	// Add the authentication cookie
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWYwMmQwMzYtYjQ1Ny00M2ExLThmYzktNWM2NDBjM2Y3ZDJhIn0.F2NM790xbzXL6b-gpxg3xUp1G76ZHS43Gy0dZwGlmJg",
		Path:  "/",
	})
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
