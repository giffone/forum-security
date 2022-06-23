package app

import (
	"github.com/giffone/forum-security/internal/adapters/api"
	"github.com/giffone/forum-security/internal/adapters/api/authentication"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/service"
)

func (a *App) authentication(auth *config.Auth,
	srvUser service.User, sMid service.Middleware, aMid api.Middleware,
) {
	authentication.NewHandler(auth, srvUser, sMid).Register(a.ctx, a.mux, aMid)
}
