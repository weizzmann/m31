package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"m31/internal/app"
	"m31/internal/config"
	"m31/pkg/client/mongodb"
	"m31/pkg/logging"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	logger := logging.GetLogger()

	logger.Info("create router")
	router := chi.NewRouter()

	logger.Info("get config")
	cfg := config.GetConfig()

	logger.Info("connection to mongodb")
	mongoDBClient, err := mongodb.NewClientCloud(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("create app")
	app, err := app.NewApp(logger, mongoDBClient, cfg.MongoDB.Collection)
	if err != nil {
		logger.Fatal(err)
	}
	app.Handler.Register(router)

	logger.Info("run app")
	go func() {
		syscall := <-c
		logger.Infof("system call:%+v", syscall)
		cancel()
	}()
	if err = app.Serve(ctx, router, "127.0.0.1:1234"); err != nil {
		logger.Infof("failed to serve:+%v\n", err)
	}
}
