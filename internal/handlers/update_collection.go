package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/compress"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"io"
	"net/http"
)

func (ctx *Handlers) UpdateCollection(res http.ResponseWriter, req *http.Request) {
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

	var collection models.Collection

	if err := json.Unmarshal(body, &collection); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	collection.Id = collectionID
	collection.UserId = userID
	err = ctx.Repos.UpdateCollection(collection)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(collection)
}
