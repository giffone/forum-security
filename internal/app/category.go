package app

import (
	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/service"
	"github.com/giffone/forum-security/internal/service/category"
)

func (a *App) category(repo repository.Repo) service.Category {
	return category.NewService(repo)
}
