package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testLogin := "testUser"
	testPassword := "testPassword"
	testUserID := "1234"

	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "Successful registration",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
				RegistrationFunc: func(id string, login string, password string) (string, error) {
					return testUserID, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User already exists",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
				GetUserFunc: func(login string) (string, error) {
					return testLogin, nil
				},
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(Reg{Login: testLogin, Password: testPassword})
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			r := httptest.NewRecorder()

			hc := &Handlers{
				Repos: tt.mockStorage,
			}
			hc.Register(r, req)

			assert.Equal(t, tt.expectedStatus, r.Code)

			if tt.expectedStatus == http.StatusOK {
				var result Result
				json.NewDecoder(r.Body).Decode(&result)
				assert.Equal(t, testUserID, result.Id)
			}
		})
	}
}
