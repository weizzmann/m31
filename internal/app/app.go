package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"m31/internal/handlers"
	"m31/internal/user"
	"m31/internal/user/db"
	"m31/pkg/logging"
)

type App struct {
	Storage user.Storage
	Service *user.Service
	Handler handlers.Handler
}

func NewApp(logger *logging.Logger, database *mongo.Database, collection string) (*App, error) {
	//logger.Info("register user storage")
	storage := db.NewStorage(logger, database, collection)
	//logger.Info("register user service")
	service := user.NewService(logger, storage)
	//logger.Info("register user handler")
	handler := user.NewHandler(logger, service)
	return &App{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}

func (a *App) Serve(ctx context.Context, chi *chi.Mux, addr string) (err error) {

	srv := &http.Server{
		Addr:           addr,
		Handler:        chi,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Service.Logger.Fatalf("listen: %+s\n", err)
		}
	}()

	a.Service.Logger.Infof("app started: %v", addr)

	<-ctx.Done()

	a.Service.Logger.Infof("app stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		a.Service.Logger.Fatalf("app shutdown Failed:%+s", err)
	}

	a.Service.Logger.Infof("app exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
