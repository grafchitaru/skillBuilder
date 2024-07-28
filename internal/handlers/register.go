package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"github.com/grafchitaru/skillBuilder/internal/users"
	"io"
	"net/http"
)

type Reg struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ctx *Handlers) Register(res http.ResponseWriter, req *http.Request) {
	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
	}

	body, ioError := io.ReadAll(reader)
	if ioError != nil {
		http.Error(res, ioError.Error(), http.StatusBadRequest)
		return
	}

	var reg Reg

	if err := json.Unmarshal(body, &reg); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	login := reg.Login
	password := reg.Password

	res.Header().Set("Content-Type", "application/json")

	_, err := ctx.Repos.GetUser(login)
	if err == nil {
		http.Error(res, "Conflict", http.StatusConflict)
		return
	}

	userID := uuid.New()

	hashedPassword, err := users.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashed password:", err)
		return
	}

	newUser, err := ctx.Repos.Registration(userID.String(), login, hashedPassword)
	if err != nil {
		fmt.Println("Error register user:", err)
		return
	}

	token, _ := auth.GenerateToken(userID, ctx.Config.SecretKey)
	auth.SetCookieAuthorization(res, req, token)

	result := models.ResultUser{
		Id:    newUser,
		Token: token,
	}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
