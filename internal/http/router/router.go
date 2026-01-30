package router

import (
	"log/slog"
	"testovoe/internal/http/handlers"
	"testovoe/internal/http/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(router *chi.Mux, h *handlers.HttpHandler, log *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscriptions", func(r chi.Router) {
			r.Post("/", h.CreateSub)
			r.Get("/", h.ListSubs)
			r.Get("/total", h.GetTotalCost)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetUserSub)
				r.Put("/", h.UpdateSub)
				r.Delete("/", h.DeleteSub)
			})
		})
	})
}
