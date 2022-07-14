package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"quest/internals/app/models"
	"time"
)

type FiatStorage struct {
	databasePool *pgxpool.Pool
}

func NewFiatStorage(pool *pgxpool.Pool) *FiatStorage {
	storage := new(FiatStorage)
	storage.databasePool = pool
	return storage

}

func (storage *FiatStorage) GetFiatLast() []models.ValCurs {
	t := time.Now()
	date := t.Format("2006-01-02")
	query := "SELECT * FROM fiat WHERE date = $1 limit 35"
	var result []models.ValCurs
	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, date)
	if err != nil {
		log.Errorln(err)
	}
	return result
}

func (storage *FiatStorage) CreateFiat(fiat models.Fiat, c chan models.Fiat) error {
	t := time.Now()
	date := t.Format("2006-01-02")
	c <- fiat
	for i := 0; i < len(fiat.Valute); i++ {

		query := "INSERT INTO fiat(date, code, value) VALUES ($1,$2, $3)"
		_, err := storage.databasePool.Exec(context.Background(), query, date, fiat.Valute[i].Code, fiat.Valute[i].Value)
		if err != nil {
			log.Errorln(err)
			return err
		}
	}
	return nil

}
func (storage *FiatStorage) GetFiatByDate(filter models.Filter) []models.ValCurs {
	args := make([]interface{}, 0)
	args = append(args, filter.Time1)

	query := "SELECT id,date,code, value from fiat WHERE date >= $1"
	if filter.Time2 != "" {
		query += "and date <= $2"
		args = append(args, filter.Time2)
	}
	var result []models.ValCurs

	err := pgxscan.Select(context.Background(), storage.databasePool, &result, query, args...)

	if err != nil {
		log.Errorln(err)
	}

	return result

}
