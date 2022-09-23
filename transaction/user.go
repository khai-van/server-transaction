package transaction

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Account string             `json:"account"`
	Balance uint64            `json:"balance"`
}

func (r *repo) GetUser(account string) (User, error) {
	var user User

	err := r.db.MongoDB.
		Collection("user").
		FindOne(context.TODO(), bson.D{{Key: "account", Value: account}}).
		Decode(&user)

	if err == mongo.ErrNoDocuments {
		return user, errors.New("User not found")
	}
	return user, err
}

func (r *repo) CreateUser(account string) (User, error) {
	if len(account) == 0 {
		return User{}, errors.New("account is empty")
	}

	doc := User{
		ID:      primitive.NewObjectID(),
		Account: account,
		Balance: 0,
	}

	res, err := r.db.MongoDB.
		Collection("user").
		InsertOne(context.TODO(), doc)

	if err != nil {
		return User{}, err
	}

	objectID := res.InsertedID.(primitive.ObjectID)
	doc.ID = objectID
	return doc, nil
}
