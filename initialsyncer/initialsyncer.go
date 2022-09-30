package initialsyncer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/event"
)

func Sync(cfg *config.Config, eventer *event.EventCreator, client client.GdriveClient) error {
	files, err := os.ReadDir(cfg.LocalDirToWatchAbsPath)

	if err != nil {
		return err
	}

	for _, f := range files {
		f := f
		go func() {
			ev := eventer.NewEvent(cfg)
			ev.SetClient(client)
			ev.SetPayload(filepath.Join(cfg.LocalDirToWatchAbsPath, f.Name()), strings.ToUpper(cfg.OnEvent))
			ev.HandleEvent()
		}()
	}

	return nil
}
