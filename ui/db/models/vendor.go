package models

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const VENDOR_ACCOUNT_TYPE = "vendor_account"

const (
	YOUTUBE_VENDOR_NAME = "YouTube"
	PCO_VENDOR_NAME     = "PCO"
)


type VendorAccount struct {
	*CommonFields    `bson:"obj_info"`
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UserId           primitive.ObjectID `bson:"user_id,omitempty"`
	Secret           string             `bson:"secret,omitempty"`
	OauthCredentials *OauthCredential   `bson:"ouath_credentials,omitempty"`
	Name             string             `bson:"name"`
}

func (va *VendorAccount) MongoId() primitive.ObjectID {
	if va.Id.IsZero() {
		now := time.Now()
		va.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return va.Id
}

func (va *VendorAccount) UpdateObjectInfo() {
	now := time.Now()
	if va.CommonFields == nil {
		va.CommonFields = new(CommonFields)
		va.EntityType = VENDOR_ACCOUNT_TYPE
		va.CreatedAt = now
	}
	va.UpdatedAt = now
}

func (va *VendorAccount) MakeRequest(req *http.Request, db *mongo.Client) error {
	if va.OauthCredentials.ExpiresAt.Before(time.Now()) {
		va.OauthCredentials.RefreshAccessToken(va.Name)
	}
}
