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
	file      string
	operation string
}

func newPayload(file, operation string) *payload {
	return &payload{
		file:      file,
		operation: operation,
	}
}

type Event struct {
	Payload *payload
	Client  client.Synchronizer
	Cfg     *config.Config
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

func (e *Event) SetPayload(file string, operation string) {
	pl := newPayload(file, operation)

	e.Payload = pl
}

func (e *Event) Handle() {
	// TODO: Add new handlers on demand
	switch {
	case e.Payload.operation == constants.Create.String():
		e.OnFileCreate()
	}
}

func (e *Event) OnFileCreate() {
	if slices.Contains(constants.AllowedOnFileCreateActions, e.Cfg.SyncAction) {
		result := e.Client.UploadFileToDir(e.Payload.file)
		e.PostProcess(result)
	} else {
		logger.Info("unable to recognize sync action")
	}
}

func (e *Event) PostProcess(result bool) {
	// TODO implement retry until and max retry
	if !result {
		logger.Error(fmt.Sprintf("event processing failed: %s - %s - %s - %s - %s",
			e.Payload.file, e.Payload.operation, e.Cfg.SyncAction, e.Cfg.LocalDirToSync, e.Cfg.GdriveClient.ParentRemoteDirID))

		time.Sleep(time.Duration(e.Cfg.SyncRetryInterval) * time.Second)

		logger.Info("prcessing queued ", e.Payload.file)

		e.Handle()
	}

	if result {
		if e.Cfg.DeleteLocalFileAfterUpload && slices.Contains(constants.AllowedDeleteLocalFileAfterUpload, e.Cfg.SyncAction) {
			e.HandleFileDeleteLocalFileAfterUpload()
		}
	}
}

func (e *Event) HandleFileDeleteLocalFileAfterUpload() {
	// TODO implement retry until and max retry
	if fileExists(e.Payload.file) {
		if err := os.Remove(e.Payload.file); err != nil {
			logger.Error(fmt.Sprintf("file deletion after upload failed: %s - %s - %s - %s - %s - %s",
				err.Error(),
				e.Payload.file,
				e.Payload.operation,
				e.Cfg.SyncAction,
				e.Cfg.LocalDirToSync,
				e.Cfg.GdriveClient.ParentRemoteDirID))

			time.Sleep(time.Duration(e.Cfg.SyncRetryInterval) * time.Second)

			logger.Info("deletion queued ", e.Payload.file)

			e.HandleFileDeleteLocalFileAfterUpload()
		}
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
