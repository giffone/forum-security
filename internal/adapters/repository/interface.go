package repository

import (
	"context"
	"database/sql"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/internal/object/dto"
	"github.com/giffone/forum-security/internal/object/model"
)

type New interface {
	Connect() *Configuration
	DataBase(ctx context.Context) *sql.DB
	Query() *object.Query
}

type Repo interface {
	Create(ctx context.Context, d dto.DTO) (int, object.Status)
	Delete(ctx context.Context, d dto.DTO) object.Status
	GetList(ctx context.Context, m model.Models) object.Status
	GetOne(ctx context.Context, m model.Model) object.Status
	ExportSettings() (*sql.DB, string, *object.Query)
}
