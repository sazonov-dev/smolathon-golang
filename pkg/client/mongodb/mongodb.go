package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCLient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port) // conn string
	}
	fmt.Println(mongoDBURL, "----------")
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDB == "" {
			authDB = "admin"
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}

	// Connect

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error %v", err)
	}

	// Ping

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to mongoDB due to error %v", err)
	}

	//
	return client.Database(database), nil
}
