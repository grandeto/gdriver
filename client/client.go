package client

import (
	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/loader"
	"github.com/grandeto/gdrive/util"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/constants"
)

type GdriveClient interface {
	Start()
	UploadFileToDir(cfg *config.ClientArguments, fname string, dirname string) bool
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Start() {
	globalFlags := loader.LoadGlobalFlags()

	handlers := loader.LoadHandlers(globalFlags)

	cli.SetHandlers(handlers)
}

func (c *Client) UploadFileToDir(cfgArgs *config.ClientArguments, fname string, dirname string) bool {
	args := []string{string(constants.Upload), string(constants.Parent), dirname, fname}

	if cfgArgs.UseServiceAccountAuth {
		args = append(args, cfgArgs.ConfigArg, util.GetDefaultConfigDir(), cfgArgs.ServiceAccountArg, cfgArgs.AuthServiceAccountFileName)
	}

	return cli.Handle(args)
}
