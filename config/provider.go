package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ujent/go-git-app/contract"
)

const appServerHostEnv = "APP_SERVER_PORT"
const gitDBConnStrEnv = "GIT_DB_CONN_STRING"
const fsTypeEnv = "FS_TYPE"
const gitRootEnv = "GIT_ROOT"
const gitRootTest = "/home/ujent/code/go-git-app/testdata"
const gitConnStr = "root:secret@/gogit"
const gitDBConnStrTest = "root:secret@/gogittest"

//Parse - get settings from the env and parse them
func Parse() (*contract.ServerSettings, error) {

	serverPort := os.Getenv(appServerHostEnv)
	if serverPort == "" {
		panic(fmt.Sprintf("%s isn't set", appServerHostEnv))
	}

	fsTypeStr := os.Getenv(fsTypeEnv)
	t, err := strconv.Atoi(fsTypeStr)
	if err != nil {
		return nil, err
	}

	fsType := contract.ToFsType(t)
	var gitConnDB string
	var rootGitPath string

	switch fsType {
	case contract.FsTypeMySQL:
		{
			gitConnDB = os.Getenv(gitDBConnStrEnv)
			if gitConnDB == "" {
				panic(fmt.Sprintf("%s isn't set", gitConnStr))
			}
		}
	case contract.FsTypeLocal:
		{
			rootGitPath = os.Getenv(gitRootEnv)
			if rootGitPath == "" {
				panic(fmt.Sprintf("%s isn't set", gitRootEnv))
			}
		}
	default:
		{
			panic(fmt.Sprintf("%s is invalid; value: %s", fsTypeEnv, fsTypeStr))
		}

	}

	return &contract.ServerSettings{Port: serverPort, GitConnStr: gitConnDB, FsType: fsType, GitRoot: rootGitPath}, nil
}

//ParseTest - returns default values for testing usage
func ParseTest() (*contract.ServerSettings, error) {
	return &contract.ServerSettings{Port: "4000", GitConnStr: gitDBConnStrTest, GitRoot: gitRootTest, FsType: contract.FsTypeLocal}, nil
}
