package app

import (
	"github.com/giffone/forum-security/internal/adapters/api"
	user2 "github.com/giffone/forum-security/internal/adapters/api/user"
	"github.com/giffone/forum-security/internal/adapters/authentication"
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/service"
	"github.com/giffone/forum-security/internal/service/user"
)

func (a *App) user(repo repository.Repo, ses api.Middleware, auth *authentication.Auth) service.User {
	srv := user.NewService(repo)
	user2.NewHandler(srv, auth).Register(a.ctx, a.router, ses)
	return srv
}
