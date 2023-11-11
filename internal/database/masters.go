package database

import (
	"context"
	"fmt"

	"github.com/smolathon/internal/models"
	"github.com/smolathon/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MasterStorageI interface {
	CreateCard(ctx context.Context, event models.Event, eventColl *mongo.Collection) (string, error)
	DeleteCard(ctx context.Context, id string) error
}

type MasterStorage struct {
	Collection *mongo.Collection
	Logger     *logging.Logger
}

func NewMasterStorage(db *DB) *MasterStorage {
	return &MasterStorage{
		Collection: db.Collections["master"],
		Logger:     db.Logger,
	}
}

func (m *MasterStorage) CreateCard(ctx context.Context, card models.Card, cardColl *mongo.Collection) (string, error) {
	//nCtx, cancnel := context.WithTimeout(ctx, 1*time.Second)
	m.Logger.Debug("create user")
	result, err := cardColl.InsertOne(ctx, card)
	m.Logger.Trace(result)
	if err != nil {
		return "", fmt.Errorf("failed to create an event due to error: %v", err)
	}

	m.Logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID) // interface type cast
	if ok {
		return oid.Hex(), nil
	}

	m.Logger.Trace(card)
	return "", fmt.Errorf("failed to convert obj id to hex. probably oid: %s", oid)
}

func (m *MasterStorage) DeleteCard(ctx context.Context, cardId string, cardColl *mongo.Collection) error {
	objectID, err := primitive.ObjectIDFromHex(cardId)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", cardId)
	}
	filter := bson.M{"_id": objectID}

	result, err := cardColl.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query")
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no objects matched")
	}
	m.Logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}
