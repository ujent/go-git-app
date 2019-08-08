package main

import (
	"log"
	"os"

	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-mysql/mysqlfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func main() {

	logger := log.New(os.Stdout, "go-git-app:", log.LstdFlags)

	settings, err := config.ParseTest()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config was successfully parsed")

	server := newServer(settings, logger)

	err = server.Start()
	if err != nil {
		panic(err)
	}

	connStr := ""
	tableName := ""
	tableName1 := ""

	fs, err := mysqlfs.New(connStr, tableName)
	if err != nil {
		log.Fatal(err)
	}

	fs1, err := mysqlfs.New(connStr, tableName1)
	if err != nil {
		log.Fatal(err)
	}

	s := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	r, err := git.Init(s, fs1)

	if err != nil {
		log.Fatal(err)
	}

	log.Print(r)
}
