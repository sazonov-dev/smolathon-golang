package database

import (
	"github.com/smolathon/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	Collections map[string]*mongo.Collection
	Logger      *logging.Logger
}

// pass collections through the config
func NewStorage(database *mongo.Database, collections []string, logger *logging.Logger) *DB {
	return &DB{
		Collections: map[string]*mongo.Collection{collections[0]: database.Collection(collections[0]),
			collections[1]: database.Collection(collections[1])},
		Logger: logger,
	}
}
