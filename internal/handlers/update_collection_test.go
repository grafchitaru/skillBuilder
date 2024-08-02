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

func TestUpdateCollection(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	collectionID := "collection1"
	mockStorage := &mocks.MockStorage{
		UpdateCollectionFunc: func(collection models.Collection) error {
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
	r.Put("/api/collections/{id}", hc.UpdateCollection)

	collection := models.Collection{
		Id:          collectionID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      testUserID,
		Name:        "Updated Collection",
		Description: "Updated Description",
		SumXp:       sql.NullInt64{Int64: 150, Valid: true},
		Xp:          sql.NullInt64{Int64: 75, Valid: true},
	}
	collectionBytes, err := json.Marshal(collection)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/collections/collection1", bytes.NewReader(collectionBytes))
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
	var updatedCollection models.Collection
	err = json.NewDecoder(rr.Body).Decode(&updatedCollection)
	require.NoError(t, err)
	assert.Equal(t, collectionID, updatedCollection.Id)
	assert.Equal(t, testUserID, updatedCollection.UserId)
	assert.Equal(t, "Updated Collection", updatedCollection.Name)
	assert.Equal(t, "Updated Description", updatedCollection.Description)
	assert.Equal(t, int64(150), updatedCollection.SumXp.Int64)
	assert.Equal(t, int64(75), updatedCollection.Xp.Int64)
}

func TestUpdateCollection_Unauthorized(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{
		UpdateCollectionFunc: func(collection models.Collection) error {
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
	r.Put("/api/collections/{id}", hc.UpdateCollection)

	collection := models.Collection{
		Id:          "collection1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      "testUserID",
		Name:        "Updated Collection",
		Description: "Updated Description",
		SumXp:       sql.NullInt64{Int64: 150, Valid: true},
		Xp:          sql.NullInt64{Int64: 75, Valid: true},
	}
	collectionBytes, err := json.Marshal(collection)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/collections/collection1", bytes.NewReader(collectionBytes))
	req.Header.Set("Content-Type", "application/json")

	require.NoError(t, err)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateCollection_BadRequest(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		UpdateCollectionFunc: func(collection models.Collection) error {
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
	r.Put("/api/collections/{id}", hc.UpdateCollection)

	req, err := http.NewRequest("PUT", "/api/collections/collection1", bytes.NewReader([]byte("invalid json")))
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

func TestUpdateCollection_InternalServerError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	mockStorage := &mocks.MockStorage{
		UpdateCollectionFunc: func(collection models.Collection) error {
			return errors.New("internal server error")
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
	r.Put("/api/collections/{id}", hc.UpdateCollection)

	collection := models.Collection{
		Id:          "collection1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      testUserID,
		Name:        "Updated Collection",
		Description: "Updated Description",
		SumXp:       sql.NullInt64{Int64: 150, Valid: true},
		Xp:          sql.NullInt64{Int64: 75, Valid: true},
	}
	collectionBytes, err := json.Marshal(collection)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/collections/collection1", bytes.NewReader(collectionBytes))
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

func TestUpdateCollection_Gzip(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "af02d036-b457-43a1-8fc9-5c640c3f7d2a"
	collectionID := "collection1"
	mockStorage := &mocks.MockStorage{
		UpdateCollectionFunc: func(collection models.Collection) error {
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
	r.Put("/api/collections/{id}", hc.UpdateCollection)

	collection := models.Collection{
		Id:          collectionID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      testUserID,
		Name:        "Updated Collection",
		Description: "Updated Description",
		SumXp:       sql.NullInt64{Int64: 150, Valid: true},
		Xp:          sql.NullInt64{Int64: 75, Valid: true},
	}
	collectionBytes, err := json.Marshal(collection)
	require.NoError(t, err)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(collectionBytes)
	gz.Close()

	req, err := http.NewRequest("PUT", "/api/collections/collection1", &buf)
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
	var updatedCollection models.Collection
	err = json.NewDecoder(rr.Body).Decode(&updatedCollection)
	require.NoError(t, err)
	assert.Equal(t, collectionID, updatedCollection.Id)
	assert.Equal(t, testUserID, updatedCollection.UserId)
	assert.Equal(t, "Updated Collection", updatedCollection.Name)
	assert.Equal(t, "Updated Description", updatedCollection.Description)
	assert.Equal(t, int64(150), updatedCollection.SumXp.Int64)
	assert.Equal(t, int64(75), updatedCollection.Xp.Int64)
}
