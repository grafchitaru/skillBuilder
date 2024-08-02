package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"net/http"
)

func (ctx *Handlers) GetMaterial(res http.ResponseWriter, req *http.Request) {
	_, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	materialID := chi.URLParam(req, "id")
	if materialID == "" {
		http.Error(res, http.StatusText(404), http.StatusNotFound)
		return
	}

	result, err := ctx.Repos.GetMaterial(materialID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
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
