package db

import (
	"context"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type VendorTokenSource struct {
	db *DB
	vendor *models.VendorAccount
}

func (db *DB) NewVendorTokenSource(vendor *models.VendorAccount) *VendorTokenSource {
	return &VendorTokenSource{db: db, vendor: vendor}
}

//Not threadsafe, please wrap in a oauth2.RefreshToken
func (ts *VendorTokenSource) Token() *oauth2.Token {
	conf := config.Config()

	//get locking collection
	col := ts.db.client.Database(conf.Mongo.LockDb).Collection(conf.Mongo.LockCol)

	//try and aquire lock
	opts := options.InsertOne()
	res, err := col.InsertOne(context.Background(), bson.M{"token_id": ts.vendor.OauthCredentials.AccessToken},opts)
	if err != nil {
		//If we didn't get the lock. Wait until whoever did refreshed the token
		if err == mongo.ErrInvalidIndexValue {
			return ts.waitForToken()
		}
		//other error return nil
		return nil
	}

	//Refresh token we have the lock

}
