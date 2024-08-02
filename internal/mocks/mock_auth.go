package mocks

import (
	"errors"
	"net/http"
)

type MockAuthService struct {
	GetUserIDFunc func(req *http.Request, secretKey string) (string, error)
}

func (mas *MockAuthService) GetUserID(req *http.Request, secretKey string) (string, error) {
	if mas.GetUserIDFunc != nil {
		return mas.GetUserIDFunc(req, secretKey)
	}
	return "", errors.New("mock error")
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}
