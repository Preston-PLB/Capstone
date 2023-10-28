package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model interface {
	Save(client *mongo.Client) error
	Delete(client *mongo.Client) error
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
	return m.Save(db.client)
}
