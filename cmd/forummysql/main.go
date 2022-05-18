package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/app"
	_ "github.com/go-sql-driver/mysql" // import mysql library
)

func main() {
	ctx := context.Background()

	db, router, port := app.NewApp(ctx).Run("mysql")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("close database: %v\n", err)
		}
	}(db)

	log.Printf("localhost%s is listening...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Printf("listening error: %v", err)
	}
}
