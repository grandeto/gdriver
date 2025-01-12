package constants

type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	// Watch Dir Events
	Create EventType = "CREATE"

	// Client actions
	UploadFileToDir string = "uploadFileToDir"

	// Handler args
	UploadArg              string = "upload" // Upload file or directory
	ParentRemoteDirFlag    string = "--parent"
	ConfigDirFlag          string = "--config"
	ServiceAccountAuthFlag string = "--service-account"
)

var (
	AllowedSyncActions                []string = []string{UploadFileToDir}
	AllowedOnFileCreateActions        []string = []string{UploadFileToDir}
	AllowedDeleteLocalFileAfterUpload []string = []string{UploadFileToDir}
)
