package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"net/http"
)

func (ctx *Handlers) DeleteCollection(res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	collectionID := chi.URLParam(req, "id")
	if collectionID == "" {
		http.Error(res, http.StatusText(404), http.StatusNotFound)
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
