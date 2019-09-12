package config

import (
	"fmt"
	"os"

	"github.com/ujent/go-git-app/contract"
)

const appServerHostEnv = "APP_SERVER_PORT"
const gitDBConnStrEnv = "GIT_DB_CONN_STRING"
const gitConnStr = "root:secret@/gogit"
const gitDBConnStrTest = "root:secret@/gogittest"

//Parse - get settings from the env and parse them
func Parse() (*contract.ServerSettings, error) {

	serverPort := os.Getenv(appServerHostEnv)
	if serverPort == "" {
		panic(fmt.Sprintf("%s isn't set", appServerHostEnv))
	}

	gitConnStr := os.Getenv(gitDBConnStrEnv)
	if gitConnStr == "" {
		panic(fmt.Sprintf("%s isn't set", gitConnStr))
	}

	return &contract.ServerSettings{Port: serverPort, GitConnStr: gitConnStr}, nil
}

//ParseTest - returns default values for testing usage
func ParseTest() (*contract.ServerSettings, error) {
	return &contract.ServerSettings{Port: "4000", GitConnStr: gitDBConnStrTest}, nil
}
