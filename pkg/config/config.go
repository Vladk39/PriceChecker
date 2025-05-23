package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type ApiExchange struct {
	Apiexchange string `env:"apiexchange"`
}

type CurrencySymbols struct {
	Symbols string `env:"symbols"`
}
type Auth struct {
	Username string `env:"authusername"`
	Password string `env:"authpassword"`
}

type DBconnection struct {
	DBconnection string `env:"dbconn"`
}

type Config struct {
	ApiExchange
	CurrencySymbols
	Auth
	DBconnection
	ServerPort string `env:"portserver"`
}

func (c *Config) GetDBconnection() *DBconnection {
	return &c.DBconnection
}

func (c *Config) GetApiExchange() *ApiExchange {
	return &c.ApiExchange
}

func (c *Config) GetCurrencySymbols() *CurrencySymbols {
	return &c.CurrencySymbols
}

func GetConfig() (*Config, error) {
	err := godotenv.Load("../envs/.env")
	if err != nil {
		return nil, errors.Wrap(err, "ошибка загрузки енв файла")
	}

	conf := &Config{}

	err = env.Parse(conf)
	if err != nil {
		return nil, errors.Wrap(err, "не может распарсить конфиг")
	}

	return conf, nil
}
