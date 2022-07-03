package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type Database struct {
	client *mongo.Client
}

func New() *Database {
	return &Database{client: nil}
}

func (d *Database) Connect(ctx context.Context) error {
	var err error

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CONNSTR"))
	if err = clientOptions.Validate(); err != nil {
		return fmt.Errorf("error parsing MongoDB connstr: %w", err)
	}

	d.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error creating connections Pool: %w", err)
	}

	return nil
}

func (d *Database) Close(ctx context.Context) error {
	return nil
}
