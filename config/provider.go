package config

import (
	"fmt"
	"os"

	"github.com/ujent/go-git-app/contract"
)

const appServerHostEnv = "APP_SERVER_PORT"

//Parse - get settings from the env and parse them
func Parse() (*contract.ServerSettings, error) {
	serverPort := os.Getenv(appServerHostEnv)

	if serverPort == "" {
		panic(fmt.Sprintf("%s isn't set", appServerHostEnv))
	}

	return &contract.ServerSettings{Port: serverPort}, nil
}

//ParseTest - returns default values for testing usage
func ParseTest() (*contract.ServerSettings, error) {
	return &contract.ServerSettings{Port: "4000"}, nil
}
