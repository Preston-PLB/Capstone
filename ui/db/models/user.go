package models

import (
	"context"
	"errors"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const USER_TYPE string = "user"

type User struct {
	*model
	mongoId      primitive.ObjectID `bson:"_id,omitempty"`
	UserId       string             `bson:"user_id,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PassowrdHash string             `bson:"password_hash,omitempty"`
}

func (user *User) Save(client *mongo.Client) error {
	conf := config.Config()

	if user.mongoId.IsZero() {
		now := time.Now()
		user.model = &model{
			EntityType: USER_TYPE,
			CreatedAt: now,
			UpdatedAt: now,
		}
		user.UserId = uuid.New().String()
		user.mongoId = primitive.NewObjectIDFromTimestamp(now)
	}

	opts := options.Update().SetUpsert(true)

	res, err := client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).UpdateOne(context.Background(), user, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("Failed to update client properly")
	}

	return nil
}

func (user *User) Delete(client *mongo.Client) error {
	conf := config.Config()
	opts := options.Delete()

	res, err := client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).DeleteOne(context.Background(), bson.M{"_id": user.mongoId}, opts)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("There was no item to delete")
	}

	return nil
}
