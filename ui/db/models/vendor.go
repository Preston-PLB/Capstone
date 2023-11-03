package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const VENDOR_ACCOUNT_TYPE = "vendor_account"

const (
	YOUTUBE_VENDOR_NAME = "YouTube"
	PCO_VENDOR_NAME     = "PCO"
)

type OauthCredential struct {
	AccessToken  string    `bson:"access_token,omitempty" json:"access_token,omitempty"`
	ExpiresIn    int       `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	ExpiresAt    time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	TokenType    string    `bson:"token_type,omitempty" json:"token_type,omitempty"`
	Scope        string    `bson:"scope,omitempty" json:"scope,omitempty"`
	RefreshToken string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
}

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
