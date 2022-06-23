package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/server"
	"github.com/giffone/forum-security/pkg/paths"
)

type App struct {
	mux *http.ServeMux
	ctx context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		ctx: ctx,
		mux: http.NewServeMux(),
	}
}

func (a *App) Run(driver string) (*sql.DB, *http.Server) {
	repo := switcher(a.ctx, driver)
	db, port, _ := repo.ExportSettings()

	home := fmt.Sprintf("%s%s", config.HomePage, port)
	tokens := config.NewSocialToken(home)

	srvMiddleware, apiMiddleware := a.middlewareService(repo)
	srvFile := a.file(repo)
	srvUser := a.user(repo, apiMiddleware, tokens)
	a.authentication(tokens, srvUser, srvMiddleware, apiMiddleware)
	srvCategory := a.category(repo)
	srvRatio := a.ratio(repo, srvMiddleware)
	srvComment := a.comment(repo, srvRatio, srvMiddleware)
	srvPost := a.post(
		repo,
		srvCategory,
		srvComment,
		srvRatio,
		srvMiddleware,
		apiMiddleware,
		srvFile,
	)
	a.home(srvPost, srvCategory, apiMiddleware)
	a.account(srvPost, srvCategory, srvComment, srvRatio, apiMiddleware)

	dir := http.Dir("internal/web/assets")
	dirHandler := http.StripPrefix("/assets/", http.FileServer(dir))
	a.mux.Handle("/assets/", dirHandler)

	// need to create paths
	paths.CreatePaths("internal/web/assets/images/post")
	paths.CreatePaths("internal/web/assets/images/avatar")
	paths.CreatePaths("cert")

	// FOR TEST ONLY
	_, _, schema := repo.ExportSettings()
	repository.NewLoremIpsum().Run(db, schema)
	// FOR TEST ONLY
	srv := server.NewServerTLS(a.mux, port)

	return db, srv
}
