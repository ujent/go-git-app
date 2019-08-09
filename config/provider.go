package config

import (
	"fmt"
	"os"

	"github.com/ujent/go-git-app/contract"
)

const appServerHostEnv = "APP_SERVER_PORT"
const gitConnStrEnv = "GIT_CONN_STRING"
const gitConnStr = "root:secret@/gogit"
const remoteGitEnv = "REMOTE_GIT_URL"

//Parse - get settings from the env and parse them
func Parse() (*contract.ServerSettings, error) {

	serverPort := os.Getenv(appServerHostEnv)
	if serverPort == "" {
		panic(fmt.Sprintf("%s isn't set", appServerHostEnv))
	}

	gitConnStr := os.Getenv(gitConnStrEnv)
	if gitConnStr == "" {
		panic(fmt.Sprintf("%s isn't set", gitConnStr))
	}

	remote := os.Getenv(remoteGitEnv)
	if gitConnStr == "" {
		panic(fmt.Sprintf("%s isn't set", remoteGitEnv))
	}

	return &contract.ServerSettings{Port: serverPort, GitConnStr: gitConnStr, GitRemote: remote}, nil
}

//ParseTest - returns default values for testing usage
func ParseTest() (*contract.ServerSettings, error) {
	return &contract.ServerSettings{Port: "4000", GitConnStr: gitConnStr, GitRemote: "http://35.239.165.218:9000"}, nil
}
