package youtube

import (
	"time"

	"google.golang.org/api/youtube/v3"
)

const (
	STATUS_PRIVATE = "private"
	STATUS_PUBLIC = "public"
)

//Inserts Broadcast into youtube
func InsertBroadcast(service *youtube.Service, title string, startTime time.Time, privacyStatus string) (*youtube.LiveBroadcast, error) {
	liveBroadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              title,
			ScheduledStartTime: startTime.Format(time.RFC3339),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: privacyStatus,
		},
	}
	return service.LiveBroadcasts.Insert([]string{"snippet", "status"}, liveBroadcast).Do()
}
