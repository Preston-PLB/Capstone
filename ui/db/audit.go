package db

import (
	"context"
	"errors"
	"sync"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// return audit trail for user
func (db *DB) FindAuditTrailForUser(userId primitive.ObjectID) ([]models.EventRecieved, []models.ActionTaken, error) {
	conf := config.Config()

	//Build sync things
	wg := new(sync.WaitGroup)
	wg.Add(2)
	errs := make([]error, 2)

	events := []models.EventRecieved{}
	actions := []models.ActionTaken{}

	//Spawn event recieved goroutine
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		opts := options.Find()
		res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.EVENT_RECIEVED_TYPE}, opts)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			errs[0] = err
			return
		}

		err = res.All(context.Background(), &events)
		if err != nil {
			errs[0] = err
			return
		}
	}(wg)

	//Spawn action taken goroutine
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		opts := options.Find()
		res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.ACTION_MAPPING_TYPE}, opts)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return
			}
			errs[1] = err
			return
		}

		err = res.All(context.Background(), &actions)
		if err != nil {
			errs[1] = err
			return
		}
	}(wg)

	//wait for go routines to finish
	wg.Wait()

	//if there was an error return the combined error
	if err := errors.Join(errs...); err != nil {
		return nil, nil, err
	}

	return events, actions, nil
}

func (db *DB) FindEventRecievedByVendorId(id string) []models.EventRecieved {
	return []models.EventRecieved{}
}
