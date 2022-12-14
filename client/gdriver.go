package client

import (
	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/loader"
	"github.com/grandeto/gdrive/util"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/constants"
)

type GdriveClient struct {
	cfgArgs *config.ClientArguments
}

func NewGdriveClient(cfgArgs *config.ClientArguments) *GdriveClient {
	return &GdriveClient{cfgArgs}
}

func (c *GdriveClient) Start() {
	globalFlags := loader.LoadGlobalFlags()

	handlers := loader.LoadHandlers(globalFlags)

	cli.SetHandlers(handlers)
}

func (c *GdriveClient) UploadFileToDir(fname string, dirname string) bool {
	args := []string{string(constants.Upload), string(constants.Parent), dirname, fname}

	if c.cfgArgs.UseServiceAccountAuth {
		args = append(args, c.cfgArgs.ConfigArg, util.GetDefaultConfigDir(), c.cfgArgs.ServiceAccountArg, c.cfgArgs.AuthServiceAccountFileName)
	}

	return cli.Handle(args)
}
