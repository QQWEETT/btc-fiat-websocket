package processors

import (
	"quest/internals/app/db"
	"quest/internals/app/models"
)

type FiatProcessor struct {
	storage *db.FiatStorage
}

func NewFiatProcessor(storage *db.FiatStorage) *FiatProcessor {
	processor := new(FiatProcessor)
	processor.storage = storage
	return processor
}

func (processor *FiatProcessor) CreateFiat(fiat models.Fiat) error {

	return processor.storage.CreateFiat(fiat, Chan4)
}

func (processor *FiatProcessor) ListFiat(filter models.Filter) ([]models.ValCurs, error) {
	return processor.storage.GetFiatByDate(filter), nil
}

func (processor *FiatProcessor) LastFiat() ([]models.ValCurs, error) {
	return processor.storage.GetFiatLast(), nil
}
