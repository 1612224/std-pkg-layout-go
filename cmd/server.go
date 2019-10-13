package main

import (
	"database/sql"
	"log"

	"useritem/http"
	"useritem/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// setup db connection
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	// setup repos
	userRepo := &sqlite.UserRepo{DB: db}
	itemRepo := &sqlite.ItemRepo{DB: db}

	// setup server
	server := http.NewServer(userRepo, itemRepo)
	log.Fatal(http.ListenAndServe(":8080", server))
}
