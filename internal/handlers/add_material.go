package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"io"
	"net/http"
)

func (ctx *Handlers) AddMaterial(res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
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

	var material models.NewMaterial

	if err := json.Unmarshal(body, &material); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	collection, err := ctx.Repos.GetCollection(material.CollectionID, userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if collection.UserId != userID {
		http.Error(res, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := ctx.Repos.CreateMaterial(userID, material.Name, material.Description, material.TypeId, material.Xp, material.Link)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ctx.Repos.AddMaterialToCollection(material.CollectionID, id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := models.ResultId{
		Id: id,
	}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
