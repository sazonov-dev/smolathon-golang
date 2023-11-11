package database

import (
	"context"
	"fmt"

	"github.com/smolathon/internal/models"
	"github.com/smolathon/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardStorageI interface {
	FindAll(ctx context.Context) ([]models.Event, error)          // for /events
	FindOne(ctx context.Context, id string) (models.Event, error) // for events/:event_id
}

type CardStorage struct {
	Collection *mongo.Collection
	Logger     *logging.Logger
}

func NewCardStorage(database *DB) *CardStorage {
	return &CardStorage{
		Collection: database.Collections["cards"],
		Logger:     database.Logger,
	}
}

// +
func (c *CardStorage) FindAll(ctx context.Context) ([]models.Card, error) {
	var cs []models.Card
	cursor, err := c.Collection.Find(ctx, bson.M{}) // bson.M{} - ?, check doc for Find()

	if cursor.Err() != nil {

		return cs, fmt.Errorf("failed to read from cards collection, err: %v", err)
	}
	if err := cursor.All(ctx, &cs); err != nil {

		return cs, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return cs, nil
}
