package app

import (
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/service"
	"github.com/giffone/forum-security/internal/service/ratio"
)

func (a *App) ratio(repo repository.Repo, sMid service.Middleware) service.Ratio {
	srvRatio := ratio.NewService(repo, sMid)
	return srvRatio
}
