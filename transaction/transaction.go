package transaction

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	From      string             `json:"from"`
	To        string             `json:"to"`
	Amount    uint64             `json:"amount"`
	CreatedAt int64              `json:"created_at"`
	Status    Status             `json:"status"`
}

type Status int32 //enum for status transaction

const (
	Create     Status = 0
	Success           = 1
	FailInFrom        = 2
	FailInTo          = 3
)

func (r *repo) Transfer(from string, to string, amount uint64) (Transaction, error) {
	// lock all incomming transaction with 2 account from, to
	r.wg.Lock(from)
	defer r.wg.UnLock(from)

	r.wg.Lock(to)
	defer r.wg.UnLock(to)
	//validate user
	fromUser, err := r.GetUser(from)
	if err != nil {
		return Transaction{}, errors.New("user from is not existing")
	}

	if fromUser.Balance < amount {
		return Transaction{}, errors.New("user from does not has enough balance")
	}
	fromUser.Balance -= amount

	toUser, err := r.GetUser(to)
	if err != nil {
		return Transaction{}, errors.New("user from is not existing")
	}
	toUser.Balance += amount
	// make a transaction
	doc := Transaction{
		ID:        primitive.NewObjectID(),
		From:      from,
		To:        to,
		Amount:    amount,
		CreatedAt: time.Now().Unix(),
	}
	transactionColection := r.db.MongoDB.Collection("transaction")
	userColection := r.db.MongoDB.Collection("user")
	// create transaction
	_, err = transactionColection.InsertOne(context.TODO(), doc)
	if err == mongo.ErrNoDocuments {
		return Transaction{}, errors.New("Transaction not found")
	}
	// update balance from_user
	_, err = userColection.UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: fromUser.ID}},
		bson.D{{Key: "$set", Value: bson.D{{
			Key:   "balance",
			Value: fromUser.Balance,
		}}}},
	)

	if err != nil {
		// update status fail update from_user
		transactionColection.UpdateOne(context.TODO(),
			bson.D{{Key: "_id", Value: doc.ID}},
			bson.D{{Key: "$set", Value: bson.D{{
				Key:   "status",
				Value: FailInTo,
			}}}},
		)

		return Transaction{}, errors.New("transaction fail when update balance from user")
	}
	// update balance to_user
	_, err = userColection.UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: toUser.ID}},
		bson.D{{Key: "$set", Value: bson.D{{
			Key:   "balance",
			Value: toUser.Balance,
		}}}},
	)

	if err != nil {
		// update status fail update to_user
		transactionColection.UpdateOne(context.TODO(),
			bson.D{{Key: "_id", Value: doc.ID}},
			bson.D{{Key: "$set", Value: bson.D{{
				Key:   "status",
				Value: FailInFrom,
			}}}},
		)

		return Transaction{}, errors.New("transaction fail when update balance to user")
	}
	// update status success for transaction
	_, err = transactionColection.UpdateOne(context.TODO(),
		bson.D{{Key: "_id", Value: doc.ID}},
		bson.D{{Key: "$set", Value: bson.D{{
			Key:   "status",
			Value: Success,
		}}}},
	)
	if err != nil {
		log.Fatal("Error when update status success", err)
	} else {
		doc.Status = Success
	}

	return doc, nil
}

func (r *repo) GetListTransaction(account string) ([]Transaction, error) {
	var listTransaction []Transaction
	// call query search all transaction of from, to with has account, and sort by created_at descending
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.db.MongoDB.
		Collection("transaction").
		Find(context.TODO(), bson.D{
			{
				Key: "$or",
				Value: []bson.D{
					{{
						Key:   "from",
						Value: account,
					}},
					{{
						Key:   "to",
						Value: account,
					}},
				},
			},
		},
			opts,
		)

	if err != nil {
		return listTransaction, err
	}

	for cursor.Next(context.TODO()) {
		var result Transaction
		if err := cursor.Decode(&result); err != nil {
			return listTransaction, err
		}
		listTransaction = append(listTransaction, result)
	}

	return listTransaction, err
}
