package app

import (
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/service"
	"github.com/giffone/forum-security/internal/service/filemaker"
)

func (a *App) file(repo repository.Repo) service.FileMaker {
	srv := filemaker.NewService(repo)
	return srv
}
