package event

import (
	"fmt"
	"os"
	"time"

	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/constants"
	"github.com/grandeto/gdriver/logger"
	"golang.org/x/exp/slices"
)

type Eventer interface {
	NewEvent() *Event
}

type EventHandler struct {
	cfg *config.Config
}

func NewEventHandler(cfg *config.Config) *EventHandler {
	return &EventHandler{cfg}
}

func (e *EventHandler) NewEvent() *Event {
	return NewEvent(e.cfg)
}

type payload struct {
	name      string
	operation string
}

func newPayload(name, operation string) *payload {
	return &payload{
		name:      name,
		operation: operation,
	}
}

type Event struct {
	Payload *payload
	Client  client.Synchronizer
	Cfg     *config.Config
	Result  bool
}

func NewEvent(cfg *config.Config) *Event {
	return &Event{
		Cfg: cfg,
	}
}

func (e *Event) GetConfig() *config.Config {
	return e.Cfg
}

func (e *Event) SetClient(client client.Synchronizer) {
	e.Client = client
}

func (e *Event) SetPayload(name, operation string) {
	pl := newPayload(name, operation)

	e.Payload = pl
}

func (e *Event) SetResult(result bool) {
	e.Result = result
}

func (e *Event) Handle() {
	// TODO: Add new handlers on demand
	switch {
	case e.Payload.operation == constants.Create.String():
		e.OnCreate()
	}
}

func (e *Event) OnCreate() {
	switch e.Cfg.SyncAction {
	case constants.UploadFileToDir:
		result := e.Client.UploadFileToDir(e.Payload.name, e.Cfg.RemoteDir)

		e.SetResult(result)

		e.PostProcess(constants.Create)
	default:
		logger.Info("unable to recognize sync action")
	}
}

func (e *Event) PostProcess(evType constants.EventType) {
	// TODO implement retry until and max retry
	if !e.Result {
		logger.Error(fmt.Sprintf("event processing failed: %s - %s - %s - %s - %s",
			e.Payload.name, evType.String(), e.Cfg.SyncAction, e.Cfg.LocalDirToWatchAbsPath, e.Cfg.RemoteDir))

		time.Sleep(time.Duration(e.Cfg.QueueProcessingInterval) * time.Second)

		logger.Info("prcessing queued ", e.Payload.name)

		e.Handle()
	}

	if e.Result &&
		slices.Contains(
			[]string{
				constants.UploadFile,
				constants.UploadFileToDir},
			e.Cfg.SyncAction) &&
		e.Cfg.DeleteAfterUpload {
		e.HandleFileDeleteAfterUpload(evType)
	}
}

func (e *Event) HandleFileDeleteAfterUpload(evType constants.EventType) {
	// TODO implement retry until and max retry
	if fileExists(e.Payload.name) {
		if err := os.Remove(e.Payload.name); err != nil {
			logger.Error(fmt.Sprintf("file deletion after upload failed: %s - %s - %s - %s - %s - %s",
				err.Error(),
				e.Payload.name,
				evType.String(),
				e.Cfg.SyncAction,
				e.Cfg.LocalDirToWatchAbsPath,
				e.Cfg.RemoteDir))

			time.Sleep(time.Duration(e.Cfg.QueueProcessingInterval) * time.Second)

			logger.Info("deletion queued ", e.Payload.name)

			e.HandleFileDeleteAfterUpload(evType)
		}
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
