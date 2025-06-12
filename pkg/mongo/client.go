package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func Connect(ctx context.Context, cfg Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Database)

	return &Client{
		Client: client,
		DB:     db,
	}, nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}
