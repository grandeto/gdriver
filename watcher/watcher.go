package watcher

import (
	"math"
	"os"
	"sync"
	"time"

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

func Watch(cfg *config.Config, eventer event.Creator, client client.GdriveClient) (*WatchProcessor, error) {
	watcher, err := NewWatchProcessor()

	if err != nil {
		return nil, err
	}

	// Start listening for events.
	go dedupLoop(cfg, eventer, client, watcher)

	// Add a path.
	err = watcher.WatchDir(cfg.LocalDirToWatchAbsPath)

	if err != nil {
		return nil, err
	}

	return watcher, nil
}

func dedupLoop(cfg *config.Config, eventer event.Creator, client client.GdriveClient, w *WatchProcessor) {
	var (
		// Wait for new events; each new event resets the timer.
		waitFor = 3000 * time.Millisecond

		// Keep track of the timers, as path → timer.
		mu     sync.Mutex
		timers = make(map[string]*time.Timer)

		// Callback we run.
		maybeTriggerEvent = func(e fsnotify.Event) {
			file, err := os.Stat(e.Name)

			if err == nil && file.Size() != 0 {
				go func() {
					ev := eventer.NewEvent(cfg)
					ev.SetClient(client)
					ev.SetPayload(e.Name, e.Op.String())
					ev.HandleEvent()
				}()
			}

			// Don't need to remove the timer if you don't have a lot of files.
			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors():
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			logger.Error(err)
		// Read from Events.
		case e, ok := <-w.Events():
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// Get timer.
			mu.Lock()
			t, ok := timers[e.Name]
			mu.Unlock()

			// No timer yet, so create one.
			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() { maybeTriggerEvent(e) })
				t.Stop()

				mu.Lock()
				timers[e.Name] = t
				mu.Unlock()
			}

			// Reset the timer for this path, so it will start from 100ms again.
			t.Reset(waitFor)
		}
	}
}
