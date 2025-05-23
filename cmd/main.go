package main

import (
	"context"
	"net/http"
	"testgo/pkg/api"
	"testgo/pkg/client"
	"testgo/pkg/config"
	"testgo/pkg/informer"
	"testgo/pkg/repository"
	"testgo/pkg/service"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	conf, err := config.GetConfig()
	if err != nil {
		logger.Info("Ошибка получения конфигурации.", zap.Error(err))
	}

	repo, err := repository.NewRepository(&conf.DBconnection)
	if err != nil {
		logger.Error("Не удалось инициализировать репозиторий", zap.Error(err))
	}

	informer := informer.NewInformer()

	quitCh := make(chan int)
	client := client.NewClientGorutines(conf, logger)
	service := service.NewAppService(conf, logger, client, informer, repo, quitCh)
	handler := api.NewHandler(service, conf)

	go service.RunParse()

	server := api.StartServer(handler)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-quitCh
		if err := repo.DB.Close(); err == nil {
			logger.Info("Закрытие соединения к БД успешно")
		}
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Ошибка при остановке сервера", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Ошибка запуска сервера", zap.Error(err))
	} else {
		logger.Info("Сервер завершил работу")
	}
}
