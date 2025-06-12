package server

import (
	"context"
	"net/http"
	"time"

	"news-app/pkg/config"
	"news-app/pkg/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	mongo  *mongo.Client
	http   *http.Server
	router chi.Router
}

func New(cfg *config.Config, logger *zap.Logger, mongo *mongo.Client) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", ww.Status()),
				zap.Int64("size", int64(ww.BytesWritten())),
				zap.Duration("duration", time.Since(start)),
			)
		})
	})

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	return &Server{
		cfg:    cfg,
		logger: logger,
		mongo:  mongo,
		router: router,
		http: &http.Server{
			Addr:         cfg.Server.Address,
			Handler:      router,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		},
	}
}

func (s *Server) Router() chi.Router {
	return s.router
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting HTTP server",
		zap.String("addr", s.http.Addr),
		zap.Duration("read_timeout", s.http.ReadTimeout),
		zap.Duration("write_timeout", s.http.WriteTimeout),
		zap.Duration("idle_timeout", s.http.IdleTimeout),
	)

	errChan := make(chan error, 1)

	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server failed to start", zap.Error(err))
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.http.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Server shutdown error", zap.Error(err))
			return err
		}

		s.logger.Info("Server stopped gracefully")
		return nil
	case err := <-errChan:
		s.logger.Error("Server failed", zap.Error(err))
		return err
	}
}
