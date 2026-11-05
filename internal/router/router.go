// Package router routes to resources
package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tolaniverse/streamforge/internal/producer"
)

func Routes(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Post("/events", producer.PushEvent)
	})
}
