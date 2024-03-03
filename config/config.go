package config

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"

	"github.com/grandeto/gdrive/util"
	"github.com/grandeto/gdriver/constants"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	LocalDirToWatchAbsPath  string
	SyncAction              string
	DeleteAfterUpload       bool
	QueueProcessingInterval int
	GdriveClient            *GdriveClient
}

func NewConfig() *Config {
	return &Config{
		LocalDirToWatchAbsPath:  getEnv("LOCAL_DIR_TO_WATCH_ABS_PATH"),
		SyncAction:              getSyncAction("SYNC_ACTION"),
		DeleteAfterUpload:       getEnvAsBool("DELETE_AFTER_UPLOAD"),
		QueueProcessingInterval: getEnvAsInt("QUEUE_PROCESSING_INTERVAL"),
		GdriveClient: &GdriveClient{
			UseServiceAccountAuth:      getEnvAsBool("USE_SERVICE_ACCOUNT_AUTH"),
			UploadArg:                  constants.UploadArg,
			ConfigDirFlag:              constants.ConfigDirFlag,
			ConfigDir:                  getEnvWithDefault("GDRIVE_CONFIG_DIR", util.GetDefaultConfigDir()),
			ServiceAccountAuthFlag:     constants.ServiceAccountAuthFlag,
			ServiceAccountAuthFileName: getServiceAccountAuthFileName("SERVICE_ACCOUNT_AUTH_FILE_NAME"),
			ParentRemoteDirFlag:        constants.ParentRemoteDirFlag,
			ParentRemoteDirID:          getEnv("PARENT_REMOTE_DIR_ID"),
		},
	}
}

type GdriveClient struct {
	UseServiceAccountAuth      bool
	UploadArg                  string
	ConfigDirFlag              string
	ConfigDir                  string
	ServiceAccountAuthFlag     string
	ServiceAccountAuthFileName string
	ParentRemoteDirFlag        string
	ParentRemoteDirID          string
}

func getSyncAction(key string) string {
	syncAction := getEnv(key)

	if !slices.Contains(constants.AllowedSyncActions, syncAction) {
		panic(fmt.Sprintf("not allowed %s. value must be in %#v", key, constants.AllowedSyncActions))
	}

	return syncAction
}

func getServiceAccountAuthFileName(key string) string {
	useServiceAccountAuth := getEnvAsBool("USE_SERVICE_ACCOUNT_AUTH")
	serviceAccountAuthFileName := os.Getenv(key)

	if useServiceAccountAuth && serviceAccountAuthFileName == "" {
		panic(fmt.Sprintf("%s not set", key))
	}

	return serviceAccountAuthFileName
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		panic(fmt.Sprintf("%s not set", key))
	}

	return val
}

func getEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultVal
	}

	return val
}

func getEnvAsBool(key string) bool {
	val := os.Getenv(key)

	switch {
	case val == "1" || val == "true":
		return true
	case val == "0" || val == "false":
		return false
	default:
		panic(fmt.Sprintf("%s not set", key))
	}
}

func getEnvAsInt(key string) int {
	val := os.Getenv(key)

	intVal, err := strconv.Atoi(val)

	if err != nil {
		panic(fmt.Sprintf("%s not set or incorrect", key))
	}

	return intVal
}
