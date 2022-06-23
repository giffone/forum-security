package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/giffone/forum-security/internal/app"
	_ "github.com/mattn/go-sqlite3" // Import go-lite3 library (by default)
)

func main() {
	ctx := context.Background()
	db, srv := app.NewApp(ctx).Run("sqlite3")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("close database: %v\n", err)
		}
	}(db)

	log.Printf("https://localhost%s is listening...\n", srv.Addr)

	// go func() {
	// 	mux := http.NewServeMux()
	// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		http.Redirect(w, r, constant.HomePage+srv.Addr, constant.Code307)
	// 	})
	// 	server.NewServer(mux, ":8080").ListenAndServe()
	// }()

	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("listening error: %v", err)
	}
}
