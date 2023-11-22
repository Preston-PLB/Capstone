package webhooks

import "time"

type Subscription struct {
	Id string `jsonapi:"primary,Subscription" bson:"id"`
	//attrs
	Active             bool      `jsonapi:"attr,active,omitempty" bson:"active"`
	ApplicationId      string    `jsonapi:"attr,application_id,omitempty" bson:"application_id"`
	AuthenticitySecret string    `jsonapi:"attr,authenticity_secret,omitempty" bson:"authenticity_secret"`
	CreatedAt          time.Time `jsonapi:"attr,created_at,rfc3339,omitempty" bson:"created_at"`
	UpdatedAt          time.Time `jsonapi:"attr,updated_at,rfc3339,omitempty" bson:"updated_at"`
	Name               string    `jsonapi:"attr,name,omitempty" bson:"name"`
	Url                string    `jsonapi:"attr,url,omitempty" bson:"url"`
}

type WebhookSubscription struct {
	Id string `jsonapi:"primary,WebhookSubscription" bson:"id"`
	//attrs
	Active             bool      `jsonapi:"attr,active,omitempty" bson:"active"`
	ApplicationId      string    `jsonapi:"attr,application_id,omitempty" bson:"application_id"`
	AuthenticitySecret string    `jsonapi:"attr,authenticity_secret,omitempty" bson:"authenticity_secret"`
	CreatedAt          time.Time `jsonapi:"attr,created_at,rfc3339,omitempty" bson:"created_at"`
	UpdatedAt          time.Time `jsonapi:"attr,updated_at,rfc3339,omitempty" bson:"updated_at"`
	Name               string    `jsonapi:"attr,name,omitempty" bson:"name"`
	Url                string    `jsonapi:"attr,url,omitempty" bson:"url"`
}
