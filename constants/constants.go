package constants

type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	// Dir Events
	Create EventType = "CREATE"

	// Client args
	Upload            string = "upload"
	Parent            string = "--parent"
	ConfigArg         string = "--config"
	ServiceAccountArg string = "--service-account"

	// Client actions
	UploadFile      string = "uploadFile"
	UploadFileToDir string = "uploadFileToDir"
)
