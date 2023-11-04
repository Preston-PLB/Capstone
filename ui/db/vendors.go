package db

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) FindVendorAccountByUser(userId primitive.ObjectID) ([]models.VendorAccount, error) {
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

//Make
func (db *DB) MakeRequestWithAccount(req *http.Request, va *models.VendorAccount) (*http.Response, error) {
	//make new credential and save new credentials to DB
	if va.OauthCredentials.ExpiresAt.Before(time.Now()) {
		err := va.OauthCredentials.RefreshAccessToken(va.Name)
		if err != nil {
			return nil, err
		}
		err = db.SaveModel(va)
		if err != nil {
			return nil, err
		}
	}

	client := http.Client{}
	req.Header.Add("Authorization", fmt.Sprintf("%s: %s", va.OauthCredentials.TokenType, va.OauthCredentials.AccessToken))
	return client.Do(req)
}
