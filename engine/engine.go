package engine

import (
	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/event"
	"github.com/grandeto/gdriver/initialsyncer"
	"github.com/grandeto/gdriver/watcher"
)

type Engine struct {
	syncer       initialsyncer.InitialSynchronizer
	watcher      watcher.Watcher
	eventHandler event.Eventer
	client       client.Synchronizer
}

func NewEngine(
	syncer initialsyncer.InitialSynchronizer,
	watcher watcher.Watcher,
	eventHandler event.Eventer,
	client client.Synchronizer) *Engine {
	return &Engine{
		syncer,
		watcher,
		eventHandler,
		client,
	}
}

func (e *Engine) Start() {
	e.client.Start()
}

func (e *Engine) Sync() error {
	return e.syncer.Sync()
}

func (e *Engine) Watch() error {
	return e.watcher.Watch(e.eventHandler, e.client)
}
