package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"net/http"
)

func (ctx *Handlers) MarkMaterialAsIncomplete(res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	materialID := chi.URLParam(req, "id")
	if materialID == "" {
		http.Error(res, http.StatusText(404), http.StatusNotFound)
		return
	}

	err = ctx.Repos.MarkMaterialAsNotCompleted(userID, materialID)
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
