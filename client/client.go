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
	UploadFileToDir(cfg *config.Config, fname string, dirname string) bool
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

func (c *Client) UploadFileToDir(cfg *config.Config, fname string, dirname string) bool {
	args := []string{string(constants.Upload), string(constants.Parent), dirname, fname}

	if cfg.UseServiceAccountAuth {
		args = append(args, cfg.ClientArgs.ConfigArg, util.GetDefaultConfigDir(), cfg.ClientArgs.ServiceAccountArg, cfg.ClientArgs.AuthServiceAccountFileName)
	}

	return cli.Handle(args)
}
