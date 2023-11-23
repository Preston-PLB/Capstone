package db

import (
	"context"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// using userId and event string return PCO Subscriptions saved to the DB
func (db *DB) FindPcoSubscriptionForUser(userId primitive.ObjectID, eventName string) (*models.PcoSubscription, error) {
	conf := config.Config()

	opts := options.FindOne()
	res := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).FindOne(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.PCO_SUBSCRIPTION_TYPE, "details.name": eventName}, opts)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	subscription := &models.PcoSubscription{}
	err := res.Decode(subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// Okay so learned something here. Interfaces are determined implemented for the type a method is related to.
// This function is not implemented for DB it is implemented for *DB and that is important
func (db *DB) SaveSubscriptionsForUser(userId primitive.ObjectID, subscriptions ...webhooks.Subscription) error {
	mods := make([]*models.PcoSubscription, 0, len(subscriptions))
	for _, sub := range subscriptions {
		mods = append(mods, &models.PcoSubscription{
			UserId:  userId,
			Details: &sub,
		})
	}

	return saveModels(db, mods...)
}
