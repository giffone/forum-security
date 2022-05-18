package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/giffone/forum-security/internal/app"
	_ "github.com/mattn/go-sqlite3" // Import go-lite3 library (by default)
)

func main() {
	ctx := context.Background()

	db, router, port := app.NewApp(ctx).Run("sqlite3")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("close database: %v\n", err)
		}
	}(db)

	log.Printf("localhost%s is listening...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("listening error: %v", err)
	}
}
