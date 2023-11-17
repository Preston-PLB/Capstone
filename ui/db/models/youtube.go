package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/youtube/v3"
)

const YOUTUBE_BROADCAST_TYPE = "youtube_broadcast"

type YoutubeBroadcast struct {
	*CommonFields `bson:"obj_info"`
	Id            primitive.ObjectID     `bson:"_id,omitempty"`
	UserId        primitive.ObjectID     `bson:"user_id,omitempty"`
	Details       *youtube.LiveBroadcast `bson:"details,omitempty"`
}

func (obj *YoutubeBroadcast) MongoId() primitive.ObjectID {
	if obj.Id.IsZero() {
		now := time.Now()
		obj.Id = primitive.NewObjectIDFromTimestamp(now)
	}

	return obj.Id
}

func (obj *YoutubeBroadcast) UpdateObjectInfo() {
	now := time.Now()
	if obj.CommonFields == nil {
		obj.CommonFields = new(CommonFields)
		obj.EntityType = YOUTUBE_BROADCAST_TYPE
		obj.CreatedAt = now
	}
	obj.UpdatedAt = now
}
