package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafchitaru/skillBuilder/internal/users"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := mocks.NewMockStorage()
	testLogin := "testUser"
	testPassword := "testPassword"
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"

	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "User not found",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
				GetUserFunc: func(login string) (string, error) {
					return "", errors.New("user not found")
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Incorrect password",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
				GetUserFunc: func(login string) (string, error) {
					return testUserID, nil
				},
				GetUserPasswordFunc: func(login string) (string, error) {
					return "wrongPassword", nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Successful login",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
				GetUserFunc: func(login string) (string, error) {
					mockStorage.SetUserPassword(testLogin, testPassword)
					return testUserID, nil
				},
				GetUserPasswordFunc: func(login string) (string, error) {
					hashedPass, _ := users.HashPassword(testPassword)
					return hashedPass, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(Reg{Login: testLogin, Password: testPassword})
			req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			r := httptest.NewRecorder()

			hc := &Handlers{
				Config: *cfg,
				Repos:  tt.mockStorage,
			}
			hc.Login(r, req)

			fmt.Println("WTF:", tt.expectedStatus)

			assert.Equal(t, tt.expectedStatus, r.Code)

			if tt.expectedStatus == http.StatusOK {
				var result Result
				json.NewDecoder(r.Body).Decode(&result)
				assert.Equal(t, testUserID, result.Id)
			}
		})
	}
}
