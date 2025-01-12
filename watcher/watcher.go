package watcher

import (
	"math"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/event"
	"github.com/grandeto/gdriver/logger"
)

type Watcher interface {
	Watch(eventHandler event.Eventer,
		client client.Synchronizer) error
	DedupLoop(eventHandler event.Eventer, client client.Synchronizer)
}

type WatchProcessor struct {
	LocalDirToSync string
	processor      *fsnotify.Watcher
}

func NewWatchProcessor(LocalDirToSync string) (*WatchProcessor, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	return &WatchProcessor{LocalDirToSync, watcher}, nil
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

func (wp *WatchProcessor) Watch(eventHandler event.Eventer, client client.Synchronizer) error {

	// Start listening for events.
	go wp.DedupLoop(eventHandler, client)

	// Add a path.
	err := wp.WatchDir(wp.LocalDirToSync)

	if err != nil {
		return err
	}

	return nil
}

func (wp *WatchProcessor) DedupLoop(eventHandler event.Eventer, client client.Synchronizer) {
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
					ev := eventHandler.NewEvent()
					ev.SetClient(client)
					ev.SetPayload(e.Name, e.Op.String())
					ev.Handle()
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
		case err, ok := <-wp.Errors():
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			logger.Error(err)
		// Read from Events.
		case e, ok := <-wp.Events():
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
