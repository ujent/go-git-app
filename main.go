package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/config"
)

func main() {

	logger := log.New(os.Stdout, "go-git-app:", log.LstdFlags)

	settings, err := config.ParseTest()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config was successfully parsed")

	db, err := sqlx.Connect("mysql", settings.GitConnStr)

	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(db, settings, logger)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}
