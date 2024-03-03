package client

type Synchronizer interface {
	Start()
	UploadFileToDir(localFilePath string) bool
}
