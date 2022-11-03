package config

import (
	"os"
	"strconv"

	"github.com/grandeto/gdriver/constants"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	LocalDirToWatchAbsPath  string
	SyncAction              string
	RemoteDir               string
	DeleteAfterUpload       bool
	QueueProcessingInterval int
	ClientArgs              *ClientArguments
}

func NewConfig() *Config {
	return &Config{
		LocalDirToWatchAbsPath:  os.Getenv("LOCAL_DIR_TO_WATCH_ABS_PATH"),
		SyncAction:              os.Getenv("SYNC_ACTION"),
		RemoteDir:               os.Getenv("REMOTE_DIR"),
		DeleteAfterUpload:       getEnvAsBool("DELETE_AFTER_UPLOAD", true),
		QueueProcessingInterval: getEnvAsInt("QUEUE_PROCESSING_INTERVAL", 3),
		ClientArgs: &ClientArguments{
			UseServiceAccountAuth:      getEnvAsBool("USE_SERVICE_ACCOUNT_AUTH", true),
			ConfigArg:                  getEnvWithDefault("CONFIG_ARG", constants.ConfigArg),
			ServiceAccountArg:          getEnvWithDefault("AUTH_SERVICE_ACCOUNT_ARG", constants.ServiceAccountArg),
			AuthServiceAccountFileName: getAuthServiceAccountFileName("AUTH_SERVICE_ACCOUNT_FILE_NAME"),
		},
	}
}

type ClientArguments struct {
	UseServiceAccountAuth      bool
	ConfigArg                  string
	ServiceAccountArg          string
	AuthServiceAccountFileName string
}

func getEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultVal
	}

	return val
}

func getEnvAsBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)

	switch {
	case val == "1" || val == "true":
		return true
	case val == "0" || val == "false":
		return false
	default:
		return defaultVal
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	val := os.Getenv(key)

	intVal, err := strconv.Atoi(val)

	if err != nil {
		return defaultVal
	}

	return intVal
}

func getAuthServiceAccountFileName(key string) string {
	AuthServiceAccountFileName := os.Getenv("AUTH_SERVICE_ACCOUNT_FILE_NAME")

	useSericeAccountAuth := getEnvAsBool("USE_SERVICE_ACCOUNT_AUTH", true)

	if useSericeAccountAuth && AuthServiceAccountFileName == "" {
		panic("AUTH_SERVICE_ACCOUNT_FILE_NAME not set")
	}

	return AuthServiceAccountFileName
}
