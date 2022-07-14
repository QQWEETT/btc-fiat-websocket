package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"quest/internals/app/models"
	"sync"
)

type BTCFiatStorage struct {
	databasePool *pgxpool.Pool
}

func NewBTCFiatStorage(pool *pgxpool.Pool) *BTCFiatStorage {
	storage := new(BTCFiatStorage)
	storage.databasePool = pool
	return storage

}

var wg sync.WaitGroup

func (storage *BTCFiatStorage) CreateBTCFiat(fiat models.Fiat, c chan float64, c2 chan bool) error {
	btc := <-c

	if btc != 0 {
		_, err := storage.databasePool.Exec(context.Background(), "DELETE FROM btcfiat *")
		if err != nil {
			log.Errorln(err)
		}

		course := btc * fiat.Valute[10].Value

		for i := 0; i < len(fiat.Valute); i++ {
			btcFiat := course / (fiat.Valute[i].Value / float64(fiat.Valute[i].Nominal))
			query2 := "INSERT INTO btcfiat(code, value) VALUES ($1,$2)"
			_, err = storage.databasePool.Exec(context.Background(), query2, fiat.Valute[i].Code, btcFiat)
			if err != nil {
				log.Errorln(err)
				return err
			}
		}
		_, err = storage.databasePool.Exec(context.Background(), "INSERT INTO btcfiat(code, value) VALUES ($1,$2)", "RUB", course)
		if err != nil {
			log.Errorln(err)
		}
		c2 <- true

		for len(c2) > 0 {
			<-c2
		}
	}
	return nil

}

func (storage *BTCFiatStorage) GetBTCFiatLast() []models.BTCFiat {
	query := "SELECT code, value FROM btcfiat"
	var result []models.BTCFiat
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query)
	if err != nil {
		log.Errorln(err)
	}
	return result
}

func (storage *BTCFiatStorage) GetBTCFiatForSocket() []models.BTCFiat {
	var result []models.BTCFiat

	query := "SELECT code, value FROM btcfiat"
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query)
	if err != nil {
		log.Errorln(err)
	}
	return result
}
