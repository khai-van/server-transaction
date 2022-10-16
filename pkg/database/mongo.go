package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect to mongodb
func connectMongoDB(DSN string, DB string) *mongo.Database {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(DSN).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	db := client.Database(DB)

	createUserCollection(ctx, db)
	createTransactionCollection(ctx, db)

	return db
}

// create index collection "user"
func createUserCollection(ctx context.Context, db *mongo.Database) {
	boolenVar := true
	db.CreateCollection(ctx, "user", nil)

	db.Collection("user").Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"account": 1,
				},
				Options: &options.IndexOptions{Unique: &boolenVar},
			},
		},
	)
}

// create index collection "transaction"
func createTransactionCollection(ctx context.Context, db *mongo.Database) {
	db.CreateCollection(ctx, "transaction", nil)

	db.Collection("transaction").Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{},
	)
}
