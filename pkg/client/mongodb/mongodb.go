package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}
	// connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB dut to error: %v", err)
	}
	//ping
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to MongoDB dut to error: %v", err)
	}
	return client.Database(database), nil
}
