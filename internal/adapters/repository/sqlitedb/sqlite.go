package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/giffone/forum-security/internal/adapters/repository"
	"github.com/giffone/forum-security/internal/adapters/repository/sqlitedb/schema"
	"github.com/giffone/forum-security/internal/config"
	"github.com/giffone/forum-security/internal/object"
	"github.com/giffone/forum-security/pkg/paths"
	"github.com/mattn/go-sqlite3"
)

type Lite struct {
	db *sql.DB
	c  *config.DriverConf
	q  *object.Query
}

func (l *Lite) Driver() *config.DriverConf {
	l.c = config.NewSqlite()
	return l.c
}

func (l *Lite) Query() *object.Query {
	l.q = &object.Query{
		Schema: schema.Query(),
	}
	return l.q
}

func (l *Lite) DataBase(ctx context.Context) *sql.DB {
	paths.CreatePaths(config.PathDBs)
	paths.CreatePaths(config.PathDBsBackup)
	if paths.NotExist(l.c.Path) { // if main base not exist
		l.create() // create new base

		if l.BackupExist() { // have backup
			l.backup()
		}
	}
	l.open()
	l.make() // make tables without restore from backup
	repository.InsertSrc(l.db, l.q.Schema)
	return l.db
}

func (l *Lite) create() {
	log.Printf("Creating %s...\n", l.c.Path)

	file, err := os.Create(l.c.Path)
	if err != nil {
		log.Fatalf("create db: can not create base - %v\n", err)
	}
	err = file.Close()
	if err != nil {
		return
	}

	log.Printf("%s created\n", l.c.Path)
}

func (l *Lite) open() {
	var err error
	log.Printf("db connection is %s\n", l.c.Connection)
	l.db, err = sql.Open(l.c.Driver, l.c.Connection) // connection - with authentication as admin
	if err != nil {
		log.Fatalf("base: open: %v\n", err)
	}
	// hoock an connection to add custom functions into database
	conn, err := l.db.Conn(context.Background())
	if err != nil {
		log.Fatalf("base: conn: %v\n", err)
	}
	defer conn.Close()
	err = conn.Raw(func(driverConn any) error {
		sqlieConn := driverConn.(*sqlite3.SQLiteConn)
		log.Printf("sqlite authenticate is %t\n", sqlieConn.AuthEnabled())
		if err := sqlieConn.RegisterFunc("minimal_id_to_show", minimalIDToShow, true); err != nil {
			return err
		}
		if err := sqlieConn.RegisterFunc("maximum_id_to_show", maximumIDToShow, true); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("base: hook: register func: %v\n", err)
	}
}

func (l *Lite) make() {
	for _, table := range config.MakeTables() {
		l.tables(table)
	}
}

func (l *Lite) tables(table string) {
	value, ok := l.q.Schema[table]
	if !ok {
		log.Fatalf("%s was not created, no such query to make a table. Fatal exit\n", table)
	}

	val := fmt.Sprintf(value, table)
	_, err := l.db.Exec(val)
	if err != nil {
		log.Fatalf("tables: execute %s: %v. Fatal exit\n", table, err) // if can not build table in Db, stop program
	}
}
