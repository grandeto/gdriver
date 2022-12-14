package initialsyncer

import (
	"os"
	"path/filepath"

	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/constants"
	"github.com/grandeto/gdriver/event"
)

type InitialSynchronizer interface {
	Sync() error
}

type InitialSync struct {
	LocalDirToSyncAbsPath string
	eventHandler          event.Eventer
	client                client.Synchronizer
}

func NewInitialSync(
	LocalDirToSyncAbsPath string,
	eventHandler event.Eventer,
	client client.Synchronizer) *InitialSync {
	return &InitialSync{
		LocalDirToSyncAbsPath,
		eventHandler,
		client,
	}
}

func (i *InitialSync) Sync() error {
	files, err := os.ReadDir(i.LocalDirToSyncAbsPath)

	if err != nil {
		return err
	}

	for _, f := range files {
		// TODO implement current write and zero bytes skip mechanism
		f := f
		go func() {
			ev := i.eventHandler.NewEvent()
			ev.SetClient(i.client)
			ev.SetPayload(filepath.Join(i.LocalDirToSyncAbsPath, f.Name()), constants.Create.String())
			ev.Handle()
		}()
	}

	return nil
}
