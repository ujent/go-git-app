package main

import (
	"log"
	"os"

	"bitbucket.org/vishjosh/bipp-go-git/experimental-app/config"
	"bitbucket.org/vishjosh/bipp-go-git/experimental-app/contract"
	"github.com/jmoiron/sqlx"
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
