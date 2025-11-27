package repo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepository struct {
	client   *mongo.Client
	database *mongo.Database
}

// Connects to DB
func NewMongoDBRepository(uri, databaseName string) (*MongoDBRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(databaseName)

	return &MongoDBRepository{
		client:   client,
		database: database,
	}, nil
}

// Close MongoDB connection
func (r *MongoDBRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

// GetCollection 
func (r *MongoDBRepository) GetCollection(name string) *mongo.Collection {
	return r.database.Collection(name)
}
