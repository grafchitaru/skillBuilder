package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"io"
	"net/http"
)

func (ctx *Handlers) UpdateCollection(res http.ResponseWriter, req *http.Request) {
	collectionID := chi.URLParam(req, "id")
	if collectionID == "" {
		http.Error(res, "ID not found", http.StatusNotFound)
		return
	}

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

	var collection models.Collection

	if err := json.Unmarshal(body, &collection); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
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
