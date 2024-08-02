package handlers

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchCollectionMaterial(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		SearchCollectionsFunc: func(query, userID string) ([]models.Collection, error) {
			return []models.Collection{
				{
					Id:          "collection1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					UserId:      userID,
					Name:        "Test Collection 1",
					Description: "Test Description 1",
					SumXp:       sql.NullInt64{Int64: 100, Valid: true},
					Xp:          sql.NullInt64{Int64: 50, Valid: true},
				},
			}, nil
		},
		SearchMaterialsFunc: func(query string) ([]models.Material, error) {
			return []models.Material{
				{
					Id:          "material1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					UserId:      testUserID,
					Name:        "Test Material 1",
					Description: "Test Description 1",
					Xp:          50,
					Link:        "Test Link",
				},
			}, nil
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
	r.Post("/api/search", hc.SearchCollectionMaterial)

	query := TextQuery{Query: "test"}
	queryBytes, err := json.Marshal(query)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/search", bytes.NewReader(queryBytes))
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

func TestSearchCollectionMaterial_Unauthorized(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{
		SearchCollectionsFunc: func(query, userID string) ([]models.Collection, error) {
			return []models.Collection{}, nil
		},
		SearchMaterialsFunc: func(query string) ([]models.Material, error) {
			return []models.Material{}, nil
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
	r.Post("/api/search", hc.SearchCollectionMaterial)

	query := TextQuery{Query: "test"}
	queryBytes, err := json.Marshal(query)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/search", bytes.NewReader(queryBytes))
	req.Header.Set("Content-Type", "application/json")

	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestSearchCollectionMaterial_BadRequest(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		SearchCollectionsFunc: func(query, userID string) ([]models.Collection, error) {
			return []models.Collection{}, nil
		},
		SearchMaterialsFunc: func(query string) ([]models.Material, error) {
			return []models.Material{}, nil
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
	r.Post("/api/search", hc.SearchCollectionMaterial)

	req, err := http.NewRequest("POST", "/api/search", bytes.NewReader([]byte("invalid json")))
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

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSearchCollectionMaterial_Gzip(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		SearchCollectionsFunc: func(query, userID string) ([]models.Collection, error) {
			return []models.Collection{
				{
					Id:          "collection1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					UserId:      userID,
					Name:        "Test Collection 1",
					Description: "Test Description 1",
					SumXp:       sql.NullInt64{Int64: 100, Valid: true},
					Xp:          sql.NullInt64{Int64: 50, Valid: true},
				},
			}, nil
		},
		SearchMaterialsFunc: func(query string) ([]models.Material, error) {
			return []models.Material{
				{
					Id:          "material1",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					UserId:      testUserID,
					Name:        "Test Material 1",
					Description: "Test Description 1",
					Xp:          50,
					Link:        "Test Link",
				},
			}, nil
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
	r.Post("/api/search", hc.SearchCollectionMaterial)

	query := TextQuery{Query: "test"}
	queryBytes, err := json.Marshal(query)
	require.NoError(t, err)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(queryBytes)
	gz.Close()

	req, err := http.NewRequest("POST", "/api/search", &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
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
