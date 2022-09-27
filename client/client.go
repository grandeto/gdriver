package client

import (
	"github.com/grandeto/gdrive/cli"
	"github.com/grandeto/gdrive/loader"
	"github.com/grandeto/gdriver/constants"
)

type GdriveClient interface {
	Start()
	UploadFile(fname string) bool
	UploadFileToDir(fname string, dirname string) bool
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

func (c *Client) UploadFile(fname string) bool {
	return cli.Handle([]string{string(constants.Upload), fname})
}

func (c *Client) UploadFileToDir(fname string, dirname string) bool {
	return cli.Handle([]string{string(constants.Upload), string(constants.Parent), dirname, fname})
}
