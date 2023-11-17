package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const EVENT_RECIEVED_TYPE = "audit_event_recieved"

type EventRecieved struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        primitive.ObjectID `bson:"user_id,omitempty"`
	Source        primitive.ObjectID `bson:"source_id,omitempty"`
}

func (obj *EventRecieved) MongoId() primitive.ObjectID {
	if obj.Id.IsZero() {
		now := time.Now()
		obj.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return obj.Id
}

func (obj *EventRecieved) UpdateObjectInfo() {
	now := time.Now()
	if obj.CommonFields == nil {
		obj.CommonFields = new(CommonFields)
		obj.EntityType = EVENT_RECIEVED_TYPE
		obj.CreatedAt = now
	}
	obj.UpdatedAt = now
}
