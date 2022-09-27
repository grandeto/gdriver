package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	LocalDirAbsPath         string
	OnEvent                 string
	SyncAction              string
	RemoteDir               string
	DeleteAfterUpload       bool
	QueueProcessingInterval int
}

func NewConfig() *Config {
	return &Config{
		LocalDirAbsPath:         os.Getenv("LOCAL_DIR_ABS_PATH"),
		OnEvent:                 os.Getenv("ON_EVENT"),
		SyncAction:              os.Getenv("SYNC_ACTION"),
		RemoteDir:               os.Getenv("REMOTE_DIR"),
		DeleteAfterUpload:       getEnvAsBool("DELETE_AFTER_UPLOAD", true),
		QueueProcessingInterval: getEnvAsInt("QUEUE_PROCESSING_INTERVAL", 3),
	}
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
