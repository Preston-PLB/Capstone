package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const USER_TYPE string = "user"

type User struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID `bson:"_id"`
	Email         string             `bson:"email,omitempty"`
	PassowrdHash  string             `bson:"password_hash,omitempty"`
}

func (user *User) MongoId() primitive.ObjectID {

	if user.Id.IsZero() {
		now := time.Now()

		user.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return user.Id
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
