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

func (db *DB) FindEventsRecievedByUserId(userId primitive.ObjectID) ([]models.EventRecieved, error) {
	conf := config.Config()
	opts := options.Find()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.EVENT_RECIEVED_TYPE}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	events := []models.EventRecieved{}
	err = res.All(context.Background(), &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (db *DB) FindActionTakenByUserId(userId primitive.ObjectID) ([]models.ActionTaken, error) {
	conf := config.Config()
	opts := options.Find()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.ACTION_TAKEN_TYPE}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	events := []models.ActionTaken{}
	err = res.All(context.Background(), &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

type VendorEventReport struct {
	Count int    `bson:"count"`
	Name  string `bson:"_id"`
}

func (db *DB) AggregateBroadcastReport(userId primitive.ObjectID) ([]VendorEventReport, error) {
	conf := config.Config()
	opts := options.Aggregate().SetAllowDiskUse(false)

	aggregation := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "obj_info.ent", Value: models.ACTION_TAKEN_TYPE}, {Key: "result", Value: "Created Broadcast"}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Aggregate(context.Background(), aggregation, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	events := []VendorEventReport{}
	err = res.All(context.Background(), &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (db *DB) AggregateVendorEventReport(userId primitive.ObjectID) ([]VendorEventReport, error) {
	conf := config.Config()
	opts := options.Aggregate().SetAllowDiskUse(false)

	aggregation := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "obj_info.ent", Value: models.EVENT_RECIEVED_TYPE}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$vendor_name"},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				},
			},
		},
	}

	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Aggregate(context.Background(), aggregation, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	events := []VendorEventReport{}
	err = res.All(context.Background(), &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
