package client

type Synchronizer interface {
	Start()
	UploadFileToDir(fileToSync string) bool
}
