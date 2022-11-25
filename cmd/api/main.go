package main

import (
	"log"
	"os"
	"refactoring/internal/config"
	"refactoring/internal/db"
	"refactoring/internal/service"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Can't create zap logger: ", err)
	}

	cfg := config.ParseConfig(logger)

	userStorage, err := db.CreateUserStorage(cfg.Database.Source)
	if err != nil {
		logger.Sugar().Fatal("Can't create a user storage", err)
	}

	serv := service.CreateNewService(userStorage)

	h := createNewHandler(logger, serv)
	r := h.Routes()

	signaler := make(chan os.Signal)
	startServer(logger, r, cfg.URL.Host, cfg.URL.Port, signaler)

	if err := h.service.Save(); err != nil {
		logger.Sugar().Fatal("Can't save changes", err)
	}
}
