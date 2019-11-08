package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-app/contract"
)

func main() {

	logger := log.New(os.Stdout, "go-git-app:", log.LstdFlags)

	settings, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config was successfully parsed")

	var db *sqlx.DB

	if settings.FsType == contract.FsTypeMySQL {
		db, err = sqlx.Connect("mysql", settings.GitConnStr)

		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
	}

	server, err := newServer(settings, logger, db)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
