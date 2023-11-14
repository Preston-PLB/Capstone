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

//return all vendor accounts for a user
func (db *DB) FindAllVendorAccountsByUser(userId primitive.ObjectID) ([]models.VendorAccount, error) {
	conf := config.Config()

	opts := options.Find()
	res, err := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).Find(context.Background(), bson.M{"user_id": userId, "obj_info.ent": models.VENDOR_ACCOUNT_TYPE}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	vendors := []models.VendorAccount{}
	err = res.All(context.Background(), &vendors)
	if err != nil {
		return nil, err
	}

	return vendors, nil
}

//find vendor for user by name
func (db *DB) FindVendorAccountByUser(userId primitive.ObjectID, name string) (*models.VendorAccount, error) {
	conf := config.Config()

	opts := options.FindOne()
	res := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).FindOne(context.Background(), bson.M{"user_id": userId, "name": name, "obj_info.ent": models.VENDOR_ACCOUNT_TYPE}, opts)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	vendor := &models.VendorAccount{}
	err := res.Decode(vendor)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}

//find vendoraccount by its unique id
func (db *DB) FindVendorAccountById(vendorId primitive.ObjectID) (*models.VendorAccount, error) {
	conf := config.Config()

	opts := options.FindOne()
	res := db.client.Database(conf.Mongo.EntDb).Collection(conf.Mongo.EntCol).FindOne(context.Background(), bson.M{"_id": vendorId, "obj_info.ent": models.VENDOR_ACCOUNT_TYPE}, opts)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, res.Err()
	}

	vendor := &models.VendorAccount{}
	err := res.Decode(vendor)
	if err != nil {
		return nil, err
	}

	return vendor, nil
}
