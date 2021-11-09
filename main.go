package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/crypto-trade/config"
	"github.com/crypto-trade/server"
	"github.com/crypto-trade/storage"
	"github.com/crypto-trade/token"
	"github.com/sirupsen/logrus"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()
	serve(ctx)
}

func serve(ctx context.Context) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := storage.New(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	token := token.NewJWT(cfg.JWT)

	srv := server.New(logger, cfg, db, token)

	go srv.Run()

	logger.Debug("App started...")
	defer logger.Debug("Gracefully finished!")

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
