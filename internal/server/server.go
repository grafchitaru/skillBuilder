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

	r.Post("/api/collection/create", hc.CreateCollection)

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
