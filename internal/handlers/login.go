package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/users"
	"io"
	"net/http"
)

func (ctx *Handlers) Login(res http.ResponseWriter, req *http.Request) {
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

	result := Result{
		Id: userID,
	}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	userIDuuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	token, _ := auth.GenerateToken(userIDuuid, ctx.Config.SecretKey)
	auth.SetCookieAuthorization(res, req, token)

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
