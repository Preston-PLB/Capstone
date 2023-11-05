package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const ACTION_MAPPING_TYPE = "action"

type ActionMapping struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	UserId        primitive.ObjectID `bson:"user_id,omitempty"`
	SourceEvent   *Event             `bson:"source_event,omitempty"`
	Action        *Action            `bson:"action,omitempty"`
}

type Action struct {
	VendorName string            `bson:"vendor_name,omitempty"`
	Type       string            `bson:"type,omitempty"`
	Fields     map[string]string `bson:"fields,omitempty"`
}

type Event struct {
	VendorName string `bson:"vendor_name,omitempty"`
	Key        string `bson:"key,omitempty"`
}
