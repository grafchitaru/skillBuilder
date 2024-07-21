package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"net/http"
)

func (ctx *Handlers) DeleteCollection(res http.ResponseWriter, req *http.Request) {
	collectionID := chi.URLParam(req, "id")
	if collectionID == "" {
		http.Error(res, "ID not found", http.StatusNotFound)
		return
	}

	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	err = ctx.Repos.DeleteCollection(userID, collectionID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
