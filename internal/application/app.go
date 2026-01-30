package application

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"testovoe/internal/config"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	ctx    context.Context
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger, router *chi.Mux) *Application {
	srv := &http.Server{
		Addr:         cfg.HttpServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	return &Application{
		ctx:    ctx,
		cfg:    cfg,
		log:    log,
		server: srv,
	}
}

func (a *Application) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *Application) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.log.Info("Run: starting http server", "addr", a.cfg.HttpServer.Addr)

		err := a.server.ListenAndServe()
		if err != nil {
			a.log.Error("ListenAndServe: ", err)
		}
	}()

	go func() {
		wg.Wait()
		a.log.Info("Run: stopping http server")
	}()

	return nil
}

func (a *Application) Shutdown() {
	a.log.Info("Shutdown: starting http server")

	err := a.server.Shutdown(a.ctx)
	if err != nil {
		a.log.Error("Shutdown: failed to shutdown http server", "error", err)
	}
}
