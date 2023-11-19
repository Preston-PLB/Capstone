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

// Interface for any object that wants to take advantage of the DB package
type Model interface {

	//Should return the _id field of the object if it exits
	//if it is new it should generate a new objectId
	MongoId() primitive.ObjectID

	//It is expected that this will update the CommonFields part of the model
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

// Upserts
func (db *DB) SaveModel(m Model) error {
	conf := config.Config()

	m.UpdateObjectInfo()
	opts := options.Update().SetUpsert(true)
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).UpdateOne(context.Background(), bson.M{"_id": m.MongoId()}, bson.M{"$set": m}, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("Failed to save model properly")
	}

	return nil
}

func (db *DB) SaveModels(m ...Model) error {
	conf := config.Config()

	writeEntry := make([]mongo.WriteModel, len(m))

	for index, model := range m {
		entry := mongo.NewUpdateOneModel()
		entry.SetFilter(bson.M{"_id": model.MongoId})
		entry.SetUpsert(true)
		entry.SetUpdate(bson.M{"$set": model})
		model.UpdateObjectInfo()

		writeEntry[index] = entry
	}

	opts := options.BulkWrite()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).BulkWrite(context.Background(), writeEntry, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("Failed to save models properly")
	}

	return nil
}

// For allowing more varidaic like things
func saveModels[T Model](db *DB, m ...T) error {
	conf := config.Config()

	writeEntry := make([]mongo.WriteModel, len(m))

	for index, model := range m {
		entry := mongo.NewUpdateOneModel()
		entry.SetFilter(bson.M{"_id": model.MongoId})
		entry.SetUpsert(true)
		entry.SetUpdate(bson.M{"$set": model})
		model.UpdateObjectInfo()

		writeEntry[index] = entry
	}

	opts := options.BulkWrite()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).BulkWrite(context.Background(), writeEntry, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("Failed to save models properly")
	}

	return nil
}

// Doesn't upsert
func (db *DB) InsertModel(m Model) error {
	conf := config.Config()

	m.UpdateObjectInfo()
	opts := options.InsertOne()
	_, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).InsertOne(context.Background(), m, opts)
	if err != nil {
		return err
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
