package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_TYPE string = "user"

type User struct {
	*CommonFields `bson:"obj_info"`
	mongoId       primitive.ObjectID `bson:"_id,omitempty"`
	Email         string             `bson:"email,omitempty"`
	PassowrdHash  string             `bson:"password_hash,omitempty"`
}

func (user *User) MongoId() primitive.ObjectID {

	if user.mongoId.IsZero() {
		now := time.Now()
		user.mongoId = primitive.NewObjectIDFromTimestamp(now)
	}

	return user.mongoId
}

func (user *User) UpdateObjectInfo() {
	now := time.Now()
	if user.CommonFields == nil {
		user.CommonFields = new(CommonFields)
		user.EntityType = USER_TYPE
		user.CreatedAt = now
	}
	user.UpdatedAt = now
}
