package auth

import "net/http"

type AuthService interface {
	GetUserID(req *http.Request, secretKey string) (string, error)
}
