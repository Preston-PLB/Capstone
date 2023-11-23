package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EVENT_RECIEVED_TYPE = "audit_event_recieved"
	ACTION_TAKEN_TYPE   = "audit_action_taken"
)

// Event Recieved
type EventRecieved struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        primitive.ObjectID `bson:"user_id,omitempty"`     //what user is this associated too
	VendorName    string             `bson:"vendor_name,omitempty"` //Vendor name of who sent us the event
	VendorId      string             `bson:"vendor_id:omitempty"`
	Type          string             `bson:"type,omitempty"` //type of event
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

// Action Taken
type ActionTaken struct {
	*CommonFields   `bson:"obj_info"`
	Id              primitive.ObjectID   `bson:"_id,omitempty"`
	UserId          primitive.ObjectID   `bson:"user_id,omitempty"`          //what user is this associated too
	TriggeringEvent primitive.ObjectID   `bson:"triggering_event,omitempty"` //what triggered this action to be taken
	Result          []primitive.ObjectID `bson:"result,omitempty"`           //list of entities effected or created from action
	VendorName      string               `bson:"vendor_name,omitempty"`      //Vendor name that the action was taken against
}

func (obj *ActionTaken) MongoId() primitive.ObjectID {
	if obj.Id.IsZero() {
		now := time.Now()
		obj.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return obj.Id
}

func (obj *ActionTaken) UpdateObjectInfo() {
	now := time.Now()
	if obj.CommonFields == nil {
		obj.CommonFields = new(CommonFields)
		obj.EntityType = EVENT_RECIEVED_TYPE
		obj.CreatedAt = now
	}
	obj.UpdatedAt = now
}
