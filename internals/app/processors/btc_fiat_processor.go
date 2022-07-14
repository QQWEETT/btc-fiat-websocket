package processors

import (
	"errors"
	"quest/internals/app/db"
	"quest/internals/app/models"
)

type BTCFiatProcessor struct {
	storage *db.BTCFiatStorage
}

var Chan = make(chan float64)
var Chan2 = make(chan bool, 1)
var Chan3 = make(chan []models.BTCFiat)
var Chan4 = make(chan models.Fiat)



func NewBTCFiatProcessor(storage *db.BTCFiatStorage) *BTCFiatProcessor {
	processor := new(BTCFiatProcessor)
	processor.storage = storage
	return processor
}

func (processor *BTCFiatProcessor) LastBTCFiat() ([]models.BTCFiat, error) {
	return processor.storage.GetBTCFiatLast(), nil
}

func (processor *BTCFiatProcessor) CreateBTCFiat(fiat models.Fiat) error {
	if fiat.Valute == nil {
		return errors.New("value should not be empty")
	}

	return processor.storage.CreateBTCFiat(fiat, Chan, Chan2)
}

func (processor *BTCFiatProcessor) BTCFiatForSocket() ([]models.BTCFiat, error) {
	return processor.storage.GetBTCFiatForSocket(), nil
}
