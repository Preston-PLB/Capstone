package services

import "time"

type Plan struct {
	Id string `jsonapi:"primary,Plan"`
	//attrs
	CanViewOrder          bool      `jsonapi:"attr,can_view_order,omitempty"`
	CreatedAt             time.Time `jsonapi:"attr,created_at,rfc3339,omitempty"`
	Dates                 string    `jsonapi:"attr,dates,omitempty"`
	FilesExpireAt         time.Time `jsonapi:"attr,files_expire_at,rfc3339,omitempty"`
	ItemsCount            int       `jsonapi:"attr,items_count,omitempty"`
	LastTimeAt            time.Time `jsonapi:"attr,last_time_at,rfc3339,omitempty"`
	MultiDay              bool      `jsonapi:"attr,multi_day,omitempty"`
	NeededPositiionsCount int       `jsonapi:"attr,needed_positions_count,omitempty"`
	OtherTimeCount        int       `jsonapi:"attr,other_time_count,omitempty"`
	Permissions           string    `jsonapi:"attr,permissions,omitempty"`
	PlanNotesCount        int       `jsonapi:"attr,plan_notes_count,omitempty"`
	PlanPeopleCount       int       `jsonapi:"attr,plan_people_count,omitempty"`
	PlanningCenterUrl     string    `jsonapi:"attr,planning_center_url,omitempty"`
	PerfersOrderView      bool      `jsonapi:"attr,prefers_order_view,omitempty"`
	Public                bool      `jsonapi:"attr,public,omitempty"`
	Rehearsable           bool      `jsonapi:"attr,rehearsable,omitempty"`
	RehearsableTimeCount  int       `jsonapi:"attr,rehearsable_time_count,omitempty"`
	RemindersDisabled     bool      `jsonapi:"attr,reminders_disabled,omitempty"`
	SeriesTitle           string    `jsonapi:"attr,series_title,omitempty"`
	ServiceTimeCount      int       `jsonapi:"attr,service_time_count,omitempty"`
	ShortDates            string    `jsonapi:"attr,short_dates,omitempty"`
	SortDate              time.Time `jsonapi:"attr,sort_date,rfc3339,omitempty"`
	Title                 string    `jsonapi:"attr,title,omitempty"`
	TotalLength           int       `jsonapi:"attr,total_length,omitempty"`
	UpdatedAt             time.Time `jsonapi:"attr,updated_at,rfc3339,omitempty"`
	//relations
	ServiceType             *ServiceType             `jsonapi:"relation,service_type,omitempty"`
	NextPlan                *Plan                    `jsonapi:"relation,next_plan,omitempty"`
	PreviousPlan            *Plan                    `jsonapi:"relation,previous_plan,omitempty"`
	AttachmentTypes         *[]AttachmentType        `jsonapi:"relation,AttachmentTypes,omitempty"`
	Series                  *Series                  `jsonapi:"relation,series,omitempty"`
	CreatedBy               *Person                  `jsonapi:"relation,created_by,omitempty"`
	UpdatedBy               *Person                  `jsonapi:"relation,updated_by,omitempty"`
	LinkedPublishingEpisode *LinkedPublishingEpisode `jsonapi:"relation,linked_publishing_episode,omitempty"`
}
