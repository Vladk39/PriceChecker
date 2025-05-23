package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testgo/pkg/config"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Currency struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price_24h"`
	Volume    float64 `json:"volume_24h"`
	LastPrice float64 `json:"last_trade_price"`
}

type ClientGorutines struct {
	Client *http.Client
	c      *config.Config
	logger *zap.Logger
}

func NewClientGorutines(c *config.Config, logger *zap.Logger) *ClientGorutines {
	return &ClientGorutines{
		Client: &http.Client{
			Timeout: 25 * time.Second,
		},
		c:      c,
		logger: logger,
	}
}

func (cg *ClientGorutines) StartWorkerParseCurrency(ctx context.Context, ticker string, outChan chan<- Currency, errCh chan error, wg *sync.WaitGroup) {
	var CurrencyState Currency
	defer wg.Done()
	tick := time.NewTicker(60 * time.Second)

	cg.logger.Info("горутина начала парсить", zap.String("тикер", ticker))

	doWork := func() error {
		url := fmt.Sprintf("%s%s", cg.c.Apiexchange, ticker)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return errors.Wrap(err, "ошибка формирования запроса.")
		}

		resp, err := cg.Client.Do(req)
		if err != nil {
			return errors.Wrap(err, "ошибка отправки запроса")
		}

		var Curency Currency
		err = json.NewDecoder(resp.Body).Decode(&Curency)
		if err != nil {
			return errors.Wrap(err, "ошибка разбора ответа")
		}

		resp.Body.Close()

		if Curency != CurrencyState {
			CurrencyState = Curency
			outChan <- Curency
		}
		return nil
	}

	if err := doWork(); err != nil {
		errCh <- err
	}

	for {
		select {
		case <-tick.C:
			if err := doWork(); err != nil {
				errCh <- err
			}
		case <-ctx.Done():
			cg.logger.Info("горутина завершила работу")
			return
		}
	}
}
