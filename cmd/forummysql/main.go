package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/giffone/forum-security/internal/app"
	_ "github.com/go-sql-driver/mysql" // import mysql library
)

func main() {
	ctx := context.Background()

	db, srv := app.NewApp(ctx).Run("mysql")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("close database: %v\n", err)
		}
	}(db)

	log.Printf("https://localhost%s is listening...\n", srv.Addr)

	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("listening error: %v", err)
	}
}
