package post

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// RegisterRoutes sets up all routes for the post handler
func RegisterRoutes(r chi.Router, h *Handler, logger *zap.Logger) {
	// Common middleware for all routes
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Custom logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	})

	// Web routes
	r.Group(func(r chi.Router) {
		r.Get("/", h.Index)
	})

	// HTMX routes
	r.Group(func(r chi.Router) {
		r.Get("/posts/{id}", h.View)
		r.Get("/posts/{id}/edit", h.EditForm)
		r.Get("/posts/{id}/delete", h.DeleteForm)
		r.Get("/posts/new", h.CreateForm)
		r.Post("/posts", h.Create)
		r.Put("/posts/{id}", h.Update)
		r.Delete("/posts/{id}", h.Delete)
	})
}
