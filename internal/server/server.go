package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/skillBuilder/internal/handlers"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/auth"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/compress"
	"github.com/grafchitaru/skillBuilder/internal/middlewares/logger"
	"github.com/rs/cors"
	"net/http"
)

func New(ctx handlers.Handlers) {
	hc := &handlers.Handlers{
		Config: ctx.Config,
		Repos:  ctx.Repos,
	}

	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{ctx.Config.ClientServer},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)
	r.Use(logger.WithLogging)
	r.Use(compress.WithCompressionResponse)
	r.Use(auth.WithUserCookie(hc.Config.SecretKey))

	r.Post("/ping", hc.Ping)

	apiRouter := chi.NewRouter()
	apiRouter.Post("/user/register", hc.Register)
	apiRouter.Post("/user/login", hc.Login)

	apiRouter.Post("/collection", hc.CreateCollection)
	apiRouter.Put("/collection/{id}", hc.UpdateCollection)
	apiRouter.Delete("/collection/{id}", hc.DeleteCollection)
	apiRouter.Get("/collection/{id}", hc.GetCollection)
	apiRouter.Get("/collections", hc.GetCollections)
	apiRouter.Get("/collections/user", hc.GetUserCollections)

	apiRouter.Post("/collection/{id}/user", hc.AddCollectionToUser)
	apiRouter.Delete("/collection/{id}/user", hc.DeleteCollectionFromUser)

	apiRouter.Post("/material", hc.AddMaterial)
	apiRouter.Put("/material/{id}", hc.UpdateMaterial)
	apiRouter.Delete("/material/{id}", hc.DeleteMaterial)
	apiRouter.Get("/material/{id}", hc.GetMaterial)
	apiRouter.Get("/collection/{id}/materials", hc.GetMaterials)

	apiRouter.Post("/material/{id}/completed", hc.MarkMaterialAsCompleted)
	apiRouter.Post("/material/{id}/incomplete", hc.MarkMaterialAsIncomplete)

	apiRouter.Post("/search", hc.SearchCollectionMaterial)

	apiRouter.Get("/material/type", hc.GetTypeMaterials)

	r.Mount("/api", apiRouter)

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
