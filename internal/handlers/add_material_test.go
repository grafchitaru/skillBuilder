package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddMaterial(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		GetCollectionFunc: func(collectionID string, userID string) (models.Collection, error) {
			return models.Collection{UserId: userID}, nil
		},
		CreateMaterialFunc: func(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
			return "test_material_id", nil
		},
		AddMaterialToCollectionFunc: func(collectionID string, materialID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.NewMaterial{
		Name:         "Test Material",
		Description:  "Test Description",
		TypeId:       "test_type_id",
		Xp:           100,
		Link:         "http://example.com",
		CollectionID: "test_collection_id",
	})
	req, err := http.NewRequest("POST", "/api/material", bytes.NewBuffer(body))
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
	hc.AddMaterial(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
}

func TestAddMaterial_Unauthorized(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{
		GetCollectionFunc: func(collectionID string, userID string) (models.Collection, error) {
			return models.Collection{}, nil
		},
		CreateMaterialFunc: func(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
			return "test_material_id", nil
		},
		AddMaterialToCollectionFunc: func(collectionID string, materialID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return "", errors.New("unauthorized")
	}

	body, _ := json.Marshal(models.NewMaterial{
		Name:         "Test Material",
		Description:  "Test Description",
		TypeId:       "test_type_id",
		Xp:           100,
		Link:         "http://example.com",
		CollectionID: "test_collection_id",
	})
	req, err := http.NewRequest("POST", "/api/material", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.AddMaterial(r, req)

	assert.Equal(t, http.StatusUnauthorized, r.Code)
}

func TestAddMaterial_Forbidden(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		GetCollectionFunc: func(collectionID string, userID string) (models.Collection, error) {
			return models.Collection{UserId: "other_user_id"}, nil
		},
		CreateMaterialFunc: func(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
			return "test_material_id", nil
		},
		AddMaterialToCollectionFunc: func(collectionID string, materialID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.NewMaterial{
		Name:         "Test Material",
		Description:  "Test Description",
		TypeId:       "test_type_id",
		Xp:           100,
		Link:         "http://example.com",
		CollectionID: "test_collection_id",
	})
	req, err := http.NewRequest("POST", "/api/material", bytes.NewBuffer(body))
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
	hc.AddMaterial(r, req)

	assert.Equal(t, http.StatusForbidden, r.Code)
}

func TestAddMaterial_CreateError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		GetCollectionFunc: func(collectionID string, userID string) (models.Collection, error) {
			return models.Collection{UserId: userID}, nil
		},
		CreateMaterialFunc: func(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
			return "", errors.New("create material error")
		},
		AddMaterialToCollectionFunc: func(collectionID string, materialID string) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.NewMaterial{
		Name:         "Test Material",
		Description:  "Test Description",
		TypeId:       "test_type_id",
		Xp:           100,
		Link:         "http://example.com",
		CollectionID: "test_collection_id",
	})
	req, err := http.NewRequest("POST", "/api/material", bytes.NewBuffer(body))
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
	hc.AddMaterial(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
}

func TestAddMaterial_AddError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		GetCollectionFunc: func(collectionID string, userID string) (models.Collection, error) {
			return models.Collection{UserId: userID}, nil
		},
		CreateMaterialFunc: func(userID string, name string, description string, typeId string, xp int, link string) (string, error) {
			return "test_material_id", nil
		},
		AddMaterialToCollectionFunc: func(collectionID string, materialID string) error {
			return errors.New("add material to collection error")
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	body, _ := json.Marshal(models.NewMaterial{
		Name:         "Test Material",
		Description:  "Test Description",
		TypeId:       "test_type_id",
		Xp:           100,
		Link:         "http://example.com",
		CollectionID: "test_collection_id",
	})
	req, err := http.NewRequest("POST", "/api/material", bytes.NewBuffer(body))
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
	hc.AddMaterial(r, req)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
}
