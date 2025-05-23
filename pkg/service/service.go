package service

import (
	"context"
	"strings"
	"sync"
	"testgo/pkg/client"
	"testgo/pkg/config"
	"testgo/pkg/informer"
	"testgo/pkg/repository"
	"time"

	"go.uber.org/zap"
)

type AppService struct {
	repo       *repository.Repository
	c          *config.Config
	client     *client.ClientGorutines
	logger     *zap.Logger
	informer   *informer.Informer
	ctx        context.Context
	cancelFunc context.CancelFunc
	quitCh     chan int
}

func NewAppService(c *config.Config, logger *zap.Logger, client *client.ClientGorutines, informer *informer.Informer, repo *repository.Repository, quitCh chan int) *AppService {
	ctx, cancel := context.WithCancel(context.Background())
	return &AppService{
		repo:       repo,
		c:          c,
		logger:     logger,
		client:     client,
		informer:   informer,
		ctx:        ctx,
		cancelFunc: cancel,
		quitCh:     quitCh,
	}
}

func (s *AppService) RunParse() {
	outCh := make(chan client.Currency, 10)
	errCh := make(chan error, 10)
	var wg sync.WaitGroup
	tickerlist := s.c.CurrencySymbols
	arrTickers := strings.Split(tickerlist.Symbols, ",")
	wg.Add(10)
	for _, value := range arrTickers {

		go s.client.StartWorkerParseCurrency(s.ctx, value, outCh, errCh, &wg)
	}

	go func() {
		wg.Wait()
		close(outCh)
		close(errCh)
	}()

	for {
		select {
		case err := <-errCh:
			s.logger.Error("получена ошибка",
				zap.Error(err))
			return
		case v := <-outCh:
			s.informer.AddInMap(v)
			s.logger.Sugar().Infof("symbol: %s | price_24h: %v | volume_24h: %v | last_trade_price %v", v.Symbol, v.Price, v.Volume, v.LastPrice)
		case <-s.ctx.Done():
			s.logger.Info("останвока парсинга")
			return
		}
	}
}

func (s *AppService) GiveCurrencyMap() map[string]informer.CurrencyInfo {
	return s.informer.GetResultMap()
}

func (s *AppService) StopWork() {
	err := s.repo.UpdateCurrency(s.GiveCurrencyMap())
	if err != nil {
		s.logger.Error("ошибка записи валюты в бд", zap.Error(err))
	}
	s.logger.Info("Запись в Базу данных успешна")
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	s.quitCh <- 1
	defer close(s.quitCh)

}
