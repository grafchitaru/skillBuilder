package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"net/http"
)

func (ctx *Handlers) MarkMaterialAsCompleted(res http.ResponseWriter, req *http.Request) {
	materialID := chi.URLParam(req, "id")
	if materialID == "" {
		http.Error(res, "ID not found", http.StatusNotFound)
		return
	}

	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	err = ctx.Repos.MarkMaterialAsCompleted(userID, materialID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := models.ResultId{}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
