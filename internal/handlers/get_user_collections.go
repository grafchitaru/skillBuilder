package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"net/http"
)

func (ctx *Handlers) GetUserCollections(res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	result, err := ctx.Repos.GetUserCollections(userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}
