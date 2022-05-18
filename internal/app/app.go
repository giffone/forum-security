package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/giffone/forum-security/internal/adapters/authentication"
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/constant"
	"github.com/giffone/forum-security/pkg/paths"
)

type App struct {
	router *http.ServeMux
	ctx    context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		ctx:    ctx,
		router: http.NewServeMux(),
	}
}

func (a *App) Run(driver string) (*sql.DB, *http.ServeMux, string) {
	repo := switcher(a.ctx, driver)
	db, port, _ := repo.ExportSettings()

	home := fmt.Sprintf("%s%s", constant.HomePage, port)
	tokens := authentication.NewSocialToken(home)

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
	a.router.Handle("/assets/", dirHandler)

	// need to create paths
	paths.CreatePaths("internal/web/assets/images/post")
	paths.CreatePaths("internal/web/assets/images/avatar")

	// FOR TEST ONLY
	_, _, schema := repo.ExportSettings()
	repository.NewLoremIpsum().Run(db, schema)
	// FOR TEST ONLY

	return db, a.router, port
}
