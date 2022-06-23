package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/giffone/forum-security/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func InsertSrc(db *sql.DB, q map[string]string) {
	src := make(map[string][]string)
	src[config.TabLikes] = []string{"like", "dislike"}
	src[config.TabCategories] = []string{
		"animal", "city", "body part", "vehicle", "actor", "18+",
		"clothing", "sport", "drink", "sea creature", "chemical element", "fruit", "country",
	}

	que := fmt.Sprintf(q[config.QueInsert6], config.TabUsers,
		config.FieldLogin, config.FieldName, config.FieldPassword, config.FieldEmail, config.FieldRoot, config.FieldCreated)
	pass, err := bcrypt.GenerateFromPassword([]byte("12345Aa"), bcrypt.MinCost)
	if err != nil {
		log.Printf("insert source: password for admin: %v\n", err)
	}
	_, err = db.Exec(que, "admin", "admin", string(pass), "admin@mail.ru", 1, time.Now()) // root=1 for admin
	if err != nil {
		log.Printf("insert source: admin did not created: %v\n", err)
	}

	for table, source := range src {
		numberLines := 0
		que = fmt.Sprintf(q[config.QueInsert2], table, config.FieldID, config.FieldBody)
		for id, value := range source {
			_, err := db.Exec(que, id+1, value)
			if err != nil {
				log.Printf("insert source: \"%s\" did not inserted to \"%s\"\n", value, table)
				continue
			}
			numberLines++
		}
		log.Printf("insert source: %d lines added to \"%s\" table\n", numberLines, table)
	}
}
