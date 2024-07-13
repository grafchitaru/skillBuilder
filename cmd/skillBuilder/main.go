package main

import (
	"fmt"
	"github.com/grafchitaru/skillBuilder/internal/config"
	"github.com/grafchitaru/skillBuilder/internal/handlers"
	"github.com/grafchitaru/skillBuilder/internal/server"
	storage2 "github.com/grafchitaru/skillBuilder/internal/storage"
	"github.com/grafchitaru/skillBuilder/internal/storage/postgresql"
)

func main() {
	cfg := *config.NewConfig()

	var storage storage2.Repositories
	var err error

	storage, err = postgresql.New(cfg.PostgresDatabaseDsn)
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	defer storage.Close()

	server.New(handlers.Handlers{Config: cfg, Repos: storage})
}
