package webhooks

import "time"

type Subscription struct {
	Id                 string    `jsonapi:"primary,Subscription"`
	//attrs
	Active             bool      `jsonapi:"attr,active,omitempty"`
	ApplicationId      string    `jsonapi:"attr,application_id,omitempty"`
	AuthenticitySecret bool      `jsonapi:"attr,authenticity_secret,omitempty"`
	CreatedAt          time.Time `jsonapi:"attr,created_at,omitempty"`
	UpdatedAt          time.Time `jsonapi:"attr,updated_at,omitempty"`
	Name               string    `jsonapi:"attr,name,omitempty"`
	Url                string    `jsonapi:"attr,url,omitempty"`
}

type WebhookSubscription struct {
	Id                 string    `jsonapi:"primary,WebhookSubscription"`
	//attrs
	Active             bool      `jsonapi:"attr,active,omitempty"`
	ApplicationId      string    `jsonapi:"attr,application_id,omitempty"`
	AuthenticitySecret bool      `jsonapi:"attr,authenticity_secret,omitempty"`
	CreatedAt          time.Time `jsonapi:"attr,created_at,omitempty"`
	UpdatedAt          time.Time `jsonapi:"attr,updated_at,omitempty"`
	Name               string    `jsonapi:"attr,name,omitempty"`
	Url                string    `jsonapi:"attr,url,omitempty"`
}
