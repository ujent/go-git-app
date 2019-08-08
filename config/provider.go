package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ujent/go-git-app/contract"
)

const appServerHostEnv = "APP_SERVER_PORT"
const logFileNameEnv = "LOG_FILE_NAME"
const logSizeMbEnv = "LOG_SIZE_MB"

//Parse - get settings from the env and parse them
func Parse() (*contract.ServerSettings, error) {
	serverPort := os.Getenv(appServerHostEnv)

	if serverPort == "" {
		panic(fmt.Sprintf("%s isn't set", appServerHostEnv))
	}

	logFile := os.Getenv(logFileNameEnv)
	if logFile == "" {
		panic(fmt.Sprintf("%s isn't set", logFileNameEnv))
	}

	logMaxSizeStr := os.Getenv(logSizeMbEnv)
	if logMaxSizeStr == "" {
		panic(fmt.Sprintf("%s isn't set", logMaxSizeStr))
	}

	logMaxSize, err := strconv.Atoi(logMaxSizeStr)
	if err != nil {
		panic(err)
	}

	if logMaxSize <= 0 {
		panic(fmt.Sprintf("%s. Wrong number %d", logSizeMbEnv, logMaxSize))
	}
	return &contract.ServerSettings{Port: serverPort, LogConfig: &contract.LogConfig{FileName: logFile, SizeMb: logMaxSize}}, nil
}

//ParseTest - returns default values for testing usage
func ParseTest() (*contract.ServerSettings, error) {
	return &contract.ServerSettings{Port: "4000", LogConfig: &contract.LogConfig{FileName: "log", SizeMb: 5}}, nil
}
