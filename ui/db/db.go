package db

import (
	"context"
	"errors"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	MongoId() primitive.ObjectID
	UpdateObjectInfo()
}

type DB struct {
	client *mongo.Client
}

func NewClient(uri string) (*DB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &DB{client: client}, nil
}

func (db *DB) SaveModel(m Model) error {
	conf := config.Config()

	m.UpdateObjectInfo()
	opts := options.Update().SetUpsert(true)
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).UpdateOne(context.Background(), bson.M{"_id": m.MongoId()}, bson.M{"$set": m}, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("Failed to update vendor account properly")
	}

	return nil
}

func (db *DB) DeleteModel(m Model) error {
	conf := config.Config()

	opts := options.Delete()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).DeleteOne(context.Background(), bson.M{"_id": m.MongoId()}, opts)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("There was no item to delete")
	}

	return nil
}
