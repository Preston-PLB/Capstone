package models

import (
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PCO_SUBSCRIPTION_TYPE = "pco_subscription"

type PcoSubscription struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID     `bson:"_id,omitempty"`
	UserId        primitive.ObjectID     `bson:"user_id,omitempty"`
	Details       *webhooks.Subscription `bson:"details,omitempty"`
}

func (obj *PcoSubscription) MongoId() primitive.ObjectID {
	if obj.Id.IsZero() {
		now := time.Now()
		obj.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return obj.Id
}

func (obj *PcoSubscription) UpdateObjectInfo() {
	now := time.Now()
	if obj.CommonFields == nil {
		obj.CommonFields = new(CommonFields)
		obj.EntityType = PCO_SUBSCRIPTION_TYPE
		obj.CreatedAt = now
	}
	obj.UpdatedAt = now
}
