package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const VENDOR_TYPE = "vendor"
const VENDOR_ACCOUNT_TYPE = "vendor_account"

type Vendor struct {
	CommonFields `bson:"obj_info"`
	mongoId      primitive.ObjectID `bson:"_id,omitempty"`
	VendorId     string             `bson:"vendor_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	OAuthUrl     string             `bson:"oauth_url"`
}

type VendorAccount struct {
	CommonFields `bson:"obj_info"`
	mongoId      primitive.ObjectID `bson:"_id,omitempty"`
	UserId       string             `bson:"user_id,omitempty"`
	Secret       string             `bson:"secret,omitempty"`
	VendorId     string             `bson:"vendor_id,omitempty"`
}
