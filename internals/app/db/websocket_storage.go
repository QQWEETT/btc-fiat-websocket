package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"quest/internals/app/models"
)

type WebsocketStorage struct {
	databasePool *pgxpool.Pool
}

func NewWebsocketStorage(pool *pgxpool.Pool) *WebsocketStorage {
	storage := new(WebsocketStorage)
	storage.databasePool = pool
	return storage

}

func (storage *WebsocketStorage) GetForSocket(c chan bool, c2 chan []models.BTCFiat) []models.BTCFiat {
	var result []models.BTCFiat
	for {
		if <-c == true {
			query := "SELECT code, value FROM btcfiat"
			err := pgxscan.Select(context.Background(), storage.databasePool, &result, query)
			if err != nil {
				log.Errorln(err)
			}
			c2 <- result
		}
	}

}
