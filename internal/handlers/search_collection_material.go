package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/models"
	"io"
	"net/http"
)

type TextQuery struct {
	Query string `json:"query"`
}

type SearchResult struct {
	Collections []models.Collection `json:"collections"`
	Materials   []models.Material   `json:"materials"`
}

func (ctx *Handlers) SearchCollectionMaterial(res http.ResponseWriter, req *http.Request) {
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

	var query TextQuery

	if err := json.Unmarshal(body, &query); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resultCollections, err := ctx.Repos.SearchCollections(query.Query, userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	resultMaterials, err := ctx.Repos.SearchMaterials(query.Query)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	searchResult := SearchResult{
		Collections: resultCollections,
		Materials:   resultMaterials,
	}

	resultData, err := json.Marshal(searchResult)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(resultData)
}
