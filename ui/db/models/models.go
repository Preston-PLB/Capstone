package models

import "time"

type CommonFields struct {
	EntityType string    `bson:"ent,omitempty"`
	CreatedAt  time.Time `bson:"created_at,omitempty"`
	UpdatedAt  time.Time `bson:"updated_at,omitempty"`
}
