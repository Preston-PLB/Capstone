package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OauthCredential struct {
	AccessToken  string    `bson:"access_token,omitempty" json:"access_token,omitempty"`
	ExpiresIn    int       `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	ExpiresAt    time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	TokenType    string    `bson:"token_type,omitempty" json:"token_type,omitempty"`
	Scope        string    `bson:"scope,omitempty" json:"scope,omitempty"`
	RefreshToken string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
}

const TOKEN_LOCK_TYPE = "token_lock"

type TokenLock struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID `bson:"_id"`
	VendorId      primitive.ObjectID `bson:"vendor_id"`
	TokenId       string             `bson:"token_id"`
	Refreshed     bool               `bson:"refreshed"`
}

func (tl *TokenLock) MongoId() primitive.ObjectID {

	if tl.Id.IsZero() {
		now := time.Now()
		tl.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return tl.Id
}

func (tl *TokenLock) UpdateObjectInfo() {
	now := time.Now()
	if tl.CommonFields == nil {
		tl.CommonFields = new(CommonFields)
		tl.EntityType = TOKEN_LOCK_TYPE
		tl.CreatedAt = now
	}
	tl.UpdatedAt = now
}
