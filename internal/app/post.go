package app

import (
	"github.com/giffone/forum-security/internal/adapters/api"
	post2 "github.com/giffone/forum-security/internal/adapters/api/post"
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/service"
	"github.com/giffone/forum-security/internal/service/post"
)

func (a *App) post(repo repository.Repo, srvCategory service.Category, srvComment service.Comment,
	srvRatio service.Ratio, sMid service.Middleware, apiMid api.Middleware, sFile service.FileMaker,
) service.Post {
	srv := post.NewService(repo, srvCategory, srvRatio, sFile, sMid)
	post2.NewHandler(srv, srvCategory, srvComment, srvRatio, sMid).Register(a.ctx, a.mux, apiMid)
	return srv
}
