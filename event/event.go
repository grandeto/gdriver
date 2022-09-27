package event

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grandeto/gdriver/client"
	"github.com/grandeto/gdriver/config"
	"github.com/grandeto/gdriver/constants"
	"github.com/grandeto/gdriver/logger"
	"golang.org/x/exp/slices"
)

type Creator interface {
	NewEvent() *Event
}

type EventCreator struct{}

func NewEventCreator() *EventCreator {
	return &EventCreator{}
}

func (ep *EventCreator) NewEvent(cfg *config.Config) *Event {
	return NewEvent(cfg)
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
	Client  client.GdriveClient
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

func (e *Event) SetClient(client client.GdriveClient) {
	e.Client = client
}

func (e *Event) SetPayload(name, operation string) {
	pl := newPayload(name, operation)

	e.Payload = pl
}

func (e *Event) SetResult(result bool) {
	e.Result = result
}

func (e *Event) HandleEvent() {
	onEvent := strings.ToUpper(e.Cfg.OnEvent)

	// TODO: Add handlers on demand
	switch {
	case onEvent == constants.Create && e.Payload.operation == constants.Create:
		e.OnCreate()
	}
}

func (e *Event) OnCreate() {
	switch e.Cfg.SyncAction {
	case constants.UploadFileToDir:
		result := e.Client.UploadFileToDir(e.Payload.name, e.Cfg.RemoteDir)

		e.SetResult(result)

		e.PostProcess()
	default:
		logger.Info("unable to recognize sync action")
	}
}

func (e *Event) PostProcess() {
	if !e.Result {
		logger.Error(fmt.Sprintf("event processing failed: %s - %s - %s - %s - %s",
			e.Payload.name, e.Cfg.OnEvent, e.Cfg.SyncAction, e.Cfg.LocalDirAbsPath, e.Cfg.RemoteDir))

		time.Sleep(time.Duration(e.Cfg.QueueProcessingInterval) * time.Second)

		logger.Info("prcessing queued ", e.Payload.name)

		e.HandleEvent()
	}

	if e.Result &&
		slices.Contains(
			[]string{
				constants.UploadFile,
				constants.UploadFileToDir},
			e.Cfg.SyncAction) &&
		e.Cfg.DeleteAfterUpload {
		e.HandleFileDeleteAfterUpload()
	}
}

func (e *Event) HandleFileDeleteAfterUpload() {
	if err := os.Remove(e.Payload.name); err != nil {
		logger.Error(fmt.Sprintf("file delete after upload failed: %s - %s - %s - %s - %s",
			e.Payload.name, e.Cfg.OnEvent, e.Cfg.SyncAction, e.Cfg.LocalDirAbsPath, e.Cfg.RemoteDir))

		time.Sleep(time.Duration(e.Cfg.QueueProcessingInterval) * time.Second)

		logger.Info("deletion queued ", e.Payload.name)

		e.HandleFileDeleteAfterUpload()
	}
}
