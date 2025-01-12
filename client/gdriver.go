package client

import (
	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/loader"
	"github.com/grandeto/gdriver/config"
)

type GdriveClient struct {
	cfg *config.GdriveClient
}

func NewGdriveClient(cfgArgs *config.GdriveClient) *GdriveClient {
	return &GdriveClient{cfgArgs}
}

func (c *GdriveClient) Start() {
	globalFlags := loader.LoadGlobalFlags()

	handlers := loader.LoadHandlers(globalFlags)

	cli.SetHandlers(handlers)
}

func (c *GdriveClient) UploadFileToDir(fileToSync string) bool {
	args := []string{c.cfg.UploadArg, c.cfg.ParentRemoteDirFlag, c.cfg.ParentRemoteDirID, fileToSync, c.cfg.ConfigDirFlag, c.cfg.ConfigDirPath}

	if c.cfg.ServiceAccountAuth {
		args = append(args, c.cfg.ServiceAccountAuthFlag, c.cfg.ServiceAccountAuthFileName)
	}

	return cli.Handle(args)
}
