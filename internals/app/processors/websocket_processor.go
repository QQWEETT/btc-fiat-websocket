package processors

import (
	"quest/internals/app/db"
	"quest/internals/app/models"
)

type WebsocketProcessor struct {
	storage *db.WebsocketStorage
}

func NewWebsocketProcessor(storage *db.WebsocketStorage) *WebsocketProcessor {
	processor := new(WebsocketProcessor)
	processor.storage = storage
	return processor
}

func (processor *WebsocketProcessor) BTCFiatForSockett() ([]models.BTCFiat, error) {
	return processor.storage.GetForSocket(Chan2, Chan3), nil
}
