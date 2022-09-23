package database

import "go.mongodb.org/mongo-driver/mongo"

type DB struct {
	MongoDB *mongo.Database
}

func (db *DB) InitMongoDB(DSN string, DB string) {
	db.MongoDB = connectMongoDB(DSN, DB)
}
