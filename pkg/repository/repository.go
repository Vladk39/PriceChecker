package repository

import (
	"testgo/pkg/config"
	"testgo/pkg/informer"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

const UpdateCurrencySQL = `UPDATE currency
SET price_24h = $2, volume_24h = $3, last_trade_price = $4
WHERE symbol = $1;`

type Repository struct {
	DB *sqlx.DB
}

// type CurrencyDB struct {
// 	Symbol    string  `db:"symbol"`
// 	Price     float32 `db:"price_24h"`
// 	Volume    float32 `db:"volume_24h"`
// 	LastPrice float32 `db:"last_trade_price"`
// }

func ConnectDB(c *config.DBconnection) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", c.DBconnection)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения к БД")
	}
	return db, nil
}

func ApplyMigrations(DB *sqlx.DB) error {
	migrationsDir := "../migrations"
	if err := goose.Up(DB.DB, migrationsDir); err != nil {
		return errors.Wrap(err, "ошибка приминения миграции")
	}
	return nil
}
func NewRepository(c *config.DBconnection) (*Repository, error) {
	db, err := ConnectDB(c)
	if err != nil {
		return nil, err
	}

	if err := ApplyMigrations(db); err != nil {
		return nil, errors.Wrap(err, "ошибка применения миграций")
	}

	return &Repository{
		DB: db,
	}, nil
}

func (repo *Repository) UpdateCurrency(info map[string]informer.CurrencyInfo) error {
	for key, value := range info {
		_, err := repo.DB.Exec(UpdateCurrencySQL, key, value.Price_24h, value.Volume_24h, value.Last_trade_price)
		if err != nil {
			return errors.Wrap(err, "не удалось записать данные в БД")
		}
	}
	return nil
}
