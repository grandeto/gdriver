package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/event"
	"github.com/grandeto/gdriver/logger"
)

type WatchProcessor struct {
	processor *fsnotify.Watcher
}

func NewWatchProcessor() (*WatchProcessor, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	return &WatchProcessor{processor: watcher}, nil
}

func (wp *WatchProcessor) WatchDir(dirpath string) error {
	return wp.processor.Add(dirpath)
}

func (wp *WatchProcessor) Events() <-chan fsnotify.Event {
	return wp.processor.Events
}

func (wp *WatchProcessor) Errors() <-chan error {
	return wp.processor.Errors
}

func (wp *WatchProcessor) Close() error {
	return wp.processor.Close()
}

func Watch(cfg *config.Config, eventer *event.EventCreator, client client.GdriveClient) (*WatchProcessor, error) {
	watcher, err := NewWatchProcessor()

	if err != nil {
		return nil, err
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case e, ok := <-watcher.Events():
				if !ok {
					return
				}

				go func() {
					ev := eventer.NewEvent(cfg)
					ev.SetClient(client)
					ev.SetPayload(e.Name, e.Op.String())
					ev.HandleEvent()
				}()
			case err, ok := <-watcher.Errors():
				if !ok {
					return
				}

				logger.Error(err)
			}
		}
	}()

	// Add a path.
	err = watcher.WatchDir(cfg.LocalDirToWatchAbsPath)

	if err != nil {
		return nil, err
	}

	return watcher, nil
}
