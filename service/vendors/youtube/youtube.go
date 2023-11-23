package youtube

import (
	"time"

	"google.golang.org/api/youtube/v3"
)

const (
	STATUS_PRIVATE = "private"
	STATUS_PUBLIC  = "public"
	ISO_8601       = "2006-01-02T15:04:05.000Z"
)

// Inserts Broadcast into youtube
func InsertBroadcast(service *youtube.Service, title string, startTime time.Time, privacyStatus string) (*youtube.LiveBroadcast, error) {
	liveBroadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              title,
			ScheduledStartTime: startTime.Format(ISO_8601),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: privacyStatus,
		},
	}
	return service.LiveBroadcasts.Insert([]string{"snippet", "status"}, liveBroadcast).Do()
}

// given a broadcast ID update the broadcast
func UpdateBroadcast(service *youtube.Service, id, title string, startTime time.Time, privacyStatus string) (*youtube.LiveBroadcast, error) {
	liveBroadcast := &youtube.LiveBroadcast{
		Id: id,
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              title,
			ScheduledStartTime: startTime.Format(ISO_8601),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: privacyStatus,
		},
	}
	return service.LiveBroadcasts.Update([]string{"snippet", "status"}, liveBroadcast).Do()
}

func DeleteBroadcast(service *youtube.Service, id string) error {
	return service.LiveBroadcasts.Delete(id).Do()
}
