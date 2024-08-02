package handlers

import (
	"github.com/grafchitaru/skillBuilder/internal/auth"
	"github.com/grafchitaru/skillBuilder/internal/config"
	"github.com/grafchitaru/skillBuilder/internal/storage"
)

type Handlers struct {
	Config config.Config
	Repos  storage.Repositories
	Auth   auth.AuthService
}
