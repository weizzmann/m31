package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"m31/internal/app"
	"m31/pkg/client/mongodb"
	"m31/pkg/logging"
)

const (
	mongoDBCollection = "users"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	logger := logging.GetLogger()

	logger.Info("create router")
	router := chi.NewRouter()

	//logger.Info("get config")
	//cfg := config.GetConfig()

	logger.Info("connection to mongodb")
	mongoDBClient, err := mongodb.NewClientCloud(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("create app")
	myApp, err := app.NewApp(logger, mongoDBClient, mongoDBCollection)
	if err != nil {
		logger.Fatal(err)
	}
	myApp.Handler.Register(router)

	logger.Info("run app")
	go func() {
		sysCall := <-c
		logger.Infof("system call:%+v", sysCall)
		cancel()
	}()

	port := os.Getenv("APP_PORT")
	if err = myApp.Serve(ctx, router, port); err != nil {
		logger.Infof("failed to serve:+%v\n", err)
	}
}
