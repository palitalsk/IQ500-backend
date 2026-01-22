package mongo

import (
	"context"
	"fmt"
	"time"

	"main/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resource struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func New(cfg *config.DB) (*Resource, error) {
	if cfg == nil {
		return nil, fmt.Errorf("mongo: DB config is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := cfg.Connection
	if uri == "" {
		if cfg.Host == "" || cfg.Port == "" {
			return nil, fmt.Errorf("mongo: connection URI and host/port are empty")
		}

		uri = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	}

	clientOpts := options.Client().ApplyURI(uri)

	if cfg.User != "" {
		clientOpts.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo: failed to connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongo: failed to ping: %w", err)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("mongo: database name is empty")
	}

	db := client.Database(cfg.Name)

	return &Resource{
		Client: client,
		DB:     db,
	}, nil
}

func (r *Resource) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Client.Disconnect(ctx)
}
