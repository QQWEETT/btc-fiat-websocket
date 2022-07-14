package processors

import (
	"errors"
	"quest/internals/app/db"
	"quest/internals/app/models"
)

type BTCProcessor struct {
	storage *db.BTCStorage
}

func NewBTCProcessor(storage *db.BTCStorage) *BTCProcessor {
	processor := new(BTCProcessor)
	processor.storage = storage
	return processor
}

func (processor *BTCProcessor) CreateBTC(btc models.BTC) error {
	if btc.Value == 0 {
		return errors.New("value should not be empty")
	}

	return processor.storage.CreateBTC(btc, Chan)
}

func (processor *BTCProcessor) ListBTC(filter models.Filter) ([]models.BTC, error) {
	return processor.storage.GetBTCByDate(filter), nil
}

func (processor *BTCProcessor) LastBTC() ([]models.BTC, error) {
	return processor.storage.GetBTCLast(), nil
}
