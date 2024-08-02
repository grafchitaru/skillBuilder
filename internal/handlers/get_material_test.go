package handlers

import (
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

func TestGetMaterial(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	materialID := "material1"
	mockStorage := &mocks.MockStorage{
		GetMaterialFunc: func(id string) (models.Material, error) {
			return models.Material{
				Id:          materialID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UserId:      testUserID,
				Name:        "Test Material",
				Description: "Test Description",
				Xp:          50,
				Link:        "Test Link",
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
	r.Get("/api/materials/{id}", hc.GetMaterial)

	req, err := http.NewRequest("GET", "/api/materials/material1", nil)
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

func TestGetMaterial_Unauthorized(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{
		GetMaterialFunc: func(id string) (models.Material, error) {
			return models.Material{}, nil
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
	r.Get("/api/materials/{id}", hc.GetMaterial)

	req, err := http.NewRequest("GET", "/api/materials/material1", nil)
	req.Header.Set("Content-Type", "application/json")

	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetMaterial_NotFound(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		GetMaterialFunc: func(id string) (models.Material, error) {
			return models.Material{}, errors.New("material not found")
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
	r.Get("/api/materials/{id}", hc.GetMaterial)

	req, err := http.NewRequest("GET", "/api/materials/material1", nil)
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
