package client

type Synchronizer interface {
	Start()
	UploadFileToDir(fname string, dirname string) bool
}
