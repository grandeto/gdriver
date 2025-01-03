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

type GdriveClient struct {
	ServiceAccountAuth         bool
	ServiceAccountAuthFlag     string
	ServiceAccountAuthFileName string
	UploadArg                  string
	ConfigDirFlag              string
	ConfigDirPath              string
	ParentRemoteDirFlag        string
	ParentRemoteDirID          string
}

type Config struct {
	LocalDirToSync             string
	SyncAction                 string
	DeleteLocalFileAfterUpload bool
	SyncRetryInterval          int
	GdriveClient               *GdriveClient
}

func NewConfig() *Config {
	return &Config{
		LocalDirToSync:             getEnv("LOCAL_DIR_TO_SYNC"),
		SyncAction:                 getEnv("SYNC_ACTION"),
		DeleteLocalFileAfterUpload: getEnvAsBool("DELETE_LOCAL_FILE_AFTER_UPLOAD"),
		SyncRetryInterval:          getEnvAsInt("SYNC_RETRY_INTERVAL"),
		GdriveClient: &GdriveClient{
			ServiceAccountAuth:         getEnvAsBool("SERVICE_ACCOUNT_AUTH"),
			ServiceAccountAuthFlag:     constants.ServiceAccountAuthFlag,
			ServiceAccountAuthFileName: getEnvWithDefault("SERVICE_ACCOUNT_AUTH_FILE_NAME", ""),
			UploadArg:                  constants.UploadArg,
			ConfigDirFlag:              constants.ConfigDirFlag,
			ConfigDirPath:              getEnvWithDefault("GDRIVE_CONFIG_DIR", util.GetDefaultConfigDir()),
			ParentRemoteDirFlag:        constants.ParentRemoteDirFlag,
			ParentRemoteDirID:          getEnv("PARENT_REMOTE_DIR_ID"),
		},
	}
}

func (c *Config) ValidateConfig() error {
	if !slices.Contains(constants.AllowedSyncActions, c.SyncAction) {
		return fmt.Errorf("not allowed %s. value must be in %#v", c.SyncAction, constants.AllowedSyncActions)
	}
	if c.GdriveClient.ServiceAccountAuth && c.GdriveClient.ServiceAccountAuthFileName == "" {
		return fmt.Errorf("SERVICE_ACCOUNT_AUTH_FILE_NAME is required when SERVICE_ACCOUNT_AUTH is set true")
	}

	return nil
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
