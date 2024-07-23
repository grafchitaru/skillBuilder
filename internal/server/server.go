package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/handlers"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/compress"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/logger"
	"net/http"
)

func New(ctx handlers.Handlers) {
	hc := &handlers.Handlers{
		Config: ctx.Config,
		Repos:  ctx.Repos,
	}

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compress.WithCompressionResponse)
	r.Use(auth.WithUserCookie(hc.Config.SecretKey))

	r.Post("/ping", hc.Ping)

	r.Post("/api/user/register", hc.Register)
	r.Post("/api/user/login", hc.Login)

	r.Post("/api/collection", hc.CreateCollection)
	r.Put("/api/collection/{id}", hc.UpdateCollection)
	r.Delete("/api/collection/{id}", hc.DeleteCollection)
	r.Get("/api/collection/{id}", hc.GetCollection)
	r.Get("/api/collections", hc.GetCollections)
	r.Get("/api/collections/user", hc.GetUserCollections)

	r.Post("/api/collection/{id}/user", hc.AddCollectionToUser)
	r.Delete("/api/collection/{id}/user", hc.DeleteCollectionFromUser)

	r.Post("/api/material", hc.AddMaterial)
	r.Put("/api/material/{id}", hc.UpdateMaterial)
	r.Delete("/api/material/{id}", hc.DeleteMaterial)
	r.Get("/api/material/{id}", hc.GetMaterial)
	r.Get("/api/collection/{id}/materials", hc.GetMaterials)

	r.Post("/api/material/{id}/completed", hc.MarkMaterialAsCompleted)
	r.Post("/api/material/{id}/incomplete", hc.MarkMaterialAsIncomplete)

	r.Post("/api/search", hc.SearchCollectionMaterial)

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
