package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"news-app/internal/server"
	"news-app/pkg/config"
	"news-app/pkg/logger"
	"news-app/pkg/mongo"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logr, err := logger.New()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logr.Sync()

	mongoClient, err := mongo.Connect(ctx, mongo.Config{
		URI:      cfg.MongoURI,
		Database: cfg.MongoDatabase,
		Timeout:  cfg.MongoTimeout,
	})
	if err != nil {
		logr.Fatal("failed to connect mongo", zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	srv := server.New(cfg, logr, mongoClient)
	srv.Handlers()
	if err := srv.Run(ctx); err != nil {
		logr.Fatal("server exited with error", zap.Error(err))
	}
}
