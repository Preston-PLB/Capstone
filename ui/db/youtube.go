package db

import (
	"context"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) FindAllBroadcastsByCorrelationId(userId primitive.ObjectID, correlationId string) ([]models.YoutubeBroadcast, error) {
	conf := config.Config()

	opts := options.Find()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.YOUTUBE_BROADCAST_TYPE, "correlation_id": correlationId}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	vendors := []models.YoutubeBroadcast{}
	err = res.All(context.Background(), &vendors)
	if err != nil {
		return nil, err
	}

	return vendors, nil
}
