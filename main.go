package main

import (
	"log"
	"os"

	"github.com/ujent/go-git-app/config"
	"github.com/ujent/go-git-app/contract"
)

const userName = "Jack Jonson"
const userEmail = "JackJonson@gmail.com"

func main() {

	logger := log.New(os.Stdout, "go-git-app:", log.LstdFlags)

	settings, err := config.ParseTest()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config was successfully parsed")

	server, err := newServer(settings, &contract.TestUsers[0], logger)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
