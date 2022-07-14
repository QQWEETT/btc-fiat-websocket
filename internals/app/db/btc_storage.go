package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"quest/internals/app/models"
)

type BTCStorage struct {
	databasePool *pgxpool.Pool
}

func NewBTCStorage(pool *pgxpool.Pool) *BTCStorage {
	storage := new(BTCStorage)
	storage.databasePool = pool
	return storage

}

func (storage *BTCStorage) GetBTCLast() []models.BTC {

	query := "SELECT id, time, value FROM btc WHERE id = (SELECT max(id) FROM btc)"
	var result []models.BTC
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query)
	if err != nil {
		log.Errorln(err)
	}
	return result
}

func (storage *BTCStorage) CreateBTC(btc models.BTC, c chan float64) error {
	s, err := storage.databasePool.Query(context.Background(), "SELECT value FROM btc WHERE id = (SELECT max(id) FROM btc)")
	if err != nil {
		log.Errorln(err)
		return err
	}

	var b float64
	for s.Next() {
		if err := s.Scan(&b); err != nil {
			return err
		}
	}
	if b != btc.Value {
		query := "INSERT INTO btc(time, value) VALUES ($1,$2)"
		_, err = storage.databasePool.Exec(context.Background(), query, btc.Time, btc.Value)
		if err != nil {
			log.Errorln(err)
			return err
		}
		c <- btc.Value

	}
	return nil
}
func (storage *BTCStorage) GetBTCByDate(filter models.Filter) []models.BTC {
	args := make([]interface{}, 0)
	args = append(args, filter.Time1)

	var result []models.BTC

	query := "SELECT id, time, value from btc WHERE time >= $1"
	if filter.Time2 != "" {
		query += "and time <= $2"
		args = append(args, filter.Time2)
	}

	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, args...)
	if err != nil {
		log.Errorln(err)
	}

	return result

}
