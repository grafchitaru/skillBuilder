package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/compress"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"github.com/grafchitaru/skillBuilder/internal/users"
	"io"
	"net/http"
)

func (ctx *Handlers) Login(res http.ResponseWriter, req *http.Request) {
	reader, err := compress.Unzip(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
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

	hashedPwd, err := ctx.Repos.GetUserPassword(login)
	if err != nil {
		http.Error(res, "User Not Found", http.StatusNotFound)
		return
	}

	if !users.ComparePasswords(hashedPwd, []byte(password)) {
		http.Error(res, "Password is not correct", http.StatusUnauthorized)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	userID, err := ctx.Repos.GetUser(login)
	if err != nil {
		http.Error(res, "User Not Found", http.StatusNotFound)
		return
	}

	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	token, _ := auth.GenerateToken(userIDuuid, ctx.Config.SecretKey)
	auth.SetCookieAuthorization(res, req, token)

	result := models.ResultUser{
		Id:    userID,
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
