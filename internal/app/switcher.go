package app

import (
	"context"
	"log"

	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/adapters/repository/mysqldb"
	"github.com/giffone/forum-security/internal/adapters/repository/sqlitedb"
	"github.com/giffone/forum-security/internal/config"
)

func switcher(ctx context.Context, driver string) repository.Repo {
	switch driver {
	case config.MysqlDriver:
		return repository.NewRepoTx(ctx, &mysqldb.MySql{})
	case config.SqliteDriver:
		return repository.NewRepo(ctx, &sqlitedb.Lite{})
	default:
		log.Fatalf("switcher: unknow driver %s\n", driver)
	}
	return nil
}
