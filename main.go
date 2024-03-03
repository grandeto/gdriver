package main

import (
	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/engine"
	"github.com/grandeto/gdriver/event"
	"github.com/grandeto/gdriver/initialsyncer"
	"github.com/grandeto/gdriver/logger"
	"github.com/grandeto/gdriver/watcher"
)

func main() {
	cfg := config.NewConfig()

	client := client.NewGdriveClient(cfg.GdriveClient)

	eventHandler := event.NewEventHandler(cfg)

	syncer := initialsyncer.NewInitialSync(cfg.LocalDirToWatchAbsPath, eventHandler, client)

	watcher, watcherErr := watcher.NewWatchProcessor(cfg.LocalDirToWatchAbsPath)

	if watcherErr != nil {
		logger.Fatal(watcherErr)
	}

	defer watcher.Close()

	eng := engine.NewEngine(syncer, watcher, eventHandler, client)

	eng.Start()

	if syncErr := eng.Sync(); syncErr != nil {
		logger.Fatal(syncErr)
	}

	if watchErr := eng.Watch(); watchErr != nil {
		logger.Fatal(watchErr)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
