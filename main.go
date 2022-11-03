package main

import (
	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/event"
	"github.com/grandeto/gdriver/initialsyncer"
	"github.com/grandeto/gdriver/logger"
	"github.com/grandeto/gdriver/watcher"
)

func main() {
	cfg := config.NewConfig()

	client := client.NewClient()

	client.Start()

	eventCreator := event.NewEventCreator()

	err := initialsyncer.SyncCreated(cfg, eventCreator, client)

	if err != nil {
		logger.Fatal(err)
	}

	watcher, err := watcher.Watch(cfg, eventCreator, client)

	if err != nil {
		logger.Fatal(err)
	}

	defer watcher.Close()

	// Block main goroutine forever.
	<-make(chan struct{})
}
