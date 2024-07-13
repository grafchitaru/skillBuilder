package handlers

import (
	"errors"
	"github.com/grafchitaru/skillBuilder/internal/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name           string
		mockStorage    *mocks.MockStorage
		expectedStatus int
	}{
		{
			name: "Error when pinging database",
			mockStorage: &mocks.MockStorage{
				PingError: errors.New("some error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Successful operation",
			mockStorage: &mocks.MockStorage{
				PingError: nil,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/ping", nil)
			require.NoError(t, err)

			hc := &Handlers{
				Repos: tt.mockStorage,
			}
			hc.Ping(r, req)

			if status := r.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
