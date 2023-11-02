package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const VENDOR_ACCOUNT_TYPE = "vendor_account"

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
	mongoId          primitive.ObjectID `bson:"_id,omitempty"`
	UserId           primitive.ObjectID `bson:"user_id,omitempty"`
	Secret           string             `bson:"secret,omitempty"`
	VendorId         string             `bson:"vendor_id,omitempty"`
	OauthCredentials *OauthCredential   `bson:"ouath_credentials,omitempty"`
}

func (va *VendorAccount) MongoId() primitive.ObjectID {
	if va.mongoId.IsZero() {
		now := time.Now()
		va.mongoId = primitive.NewObjectIDFromTimestamp(now)
	}

	return va.mongoId
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
