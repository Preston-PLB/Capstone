package db

import (
	"context"
	"errors"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type VendorTokenSource struct {
	db     *DB
	vendor *models.VendorAccount
}

func (db *DB) NewVendorTokenSource(vendor *models.VendorAccount) *VendorTokenSource {
	return &VendorTokenSource{db: db, vendor: vendor}
}

// Not threadsafe, please wrap in a oauth2.RefreshToken
func (ts *VendorTokenSource) Token() (*oauth2.Token, error) {
	conf := config.Config()

	//get locking collection
	col := ts.db.client.Database(conf.Mongo.LockDb).Collection(conf.Mongo.LockCol)

	//Define lock
	token_lock := &models.TokenLock{
		VendorId:  ts.vendor.MongoId(),
		TokenId:   ts.vendor.OauthCredentials.RefreshToken,
		Refreshed: false,
	}
	//Don't forget to create the mongo id
	token_lock.MongoId()
	token_lock.UpdateObjectInfo()

	//try and aquire lock
	opts := options.InsertOne()
	_, err := col.InsertOne(context.Background(), token_lock, opts)
	if err != nil {
		//If we didn't get the lock. Wait until whoever did refreshed the token
		if mongo.IsDuplicateKeyError(err) {
			err = ts.waitForToken(token_lock)
			if err != nil {
				return nil, err
			}

			//get new vendorAccount
			va, err := ts.db.FindVendorAccountById(ts.vendor.MongoId())
			if err != nil {
				return nil, err
			}

			//re-assign vendor account. Let go garbage collector handle the rest
			ts.vendor = va
			return ts.vendor.Token(), nil
		}
		//other error return nil
		return nil, err
	}

	//Refresh token we have the lock
	token, err := oauth2.RefreshToken(context.Background(), conf.Vendors[ts.vendor.Name].OauthConfig(), ts.vendor.OauthCredentials.RefreshToken)
	if err != nil {
		return token, err
	}

	//update vendor
	ts.vendor.OauthCredentials.RefreshToken = token.RefreshToken

	//save vendor
	err = ts.db.SaveModel(ts.vendor)
	if err != nil {
		return nil, err
	}

	//release lock
	updateOpts := options.Update()
	_, err = col.UpdateByID(context.Background(), token_lock.MongoId(), bson.M{"$set": bson.M{"refreshed": true}}, updateOpts)
	if err != nil {
		return nil, err
	}

	return token, nil
}

//Allow us to check for kind of error at the end
var TokenWaitExpired error = errors.New("Waiting for token to refresh took too long")

//Used to extract the token lock that was updated when the change stream alerts
type tokenLockChangeEvent struct {
	TokenLock *models.TokenLock `bson:"fullDocument"`
}

func (ts *VendorTokenSource) waitForToken(tl *models.TokenLock) error {
	conf := config.Config()

	//get locking collection
	col := ts.db.client.Database(conf.Mongo.LockDb).Collection(conf.Mongo.LockCol)

	//Define timeoutfunction
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	//Get change stream that looks for our token lock
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := col.Watch(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "fullDocument.vendor_id", Value: tl.VendorId}, {Key: "fullDocument.token_id", Value: tl.TokenId}}}},
	}, opts)
	if err != nil {

		return  err
	}
	defer changeStream.Close(context.Background())

	for changeStream.Next(ctx) {
		var tl_event tokenLockChangeEvent

		err := changeStream.Decode(&tl_event)
		if err != nil {
			return err
		}

		if tl_event.TokenLock.Refreshed {
			return nil
		}

	}
	return TokenWaitExpired
}
