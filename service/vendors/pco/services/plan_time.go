package services

import "time"

type PlanTime struct {
	//id
	Id string `jsonapi:"primary,PlanTime"`
	//attributes
	CreatedAt     time.Time     `jsonapi:"attr,created_at,rfc3339,omitempty"`
	StartsAt      time.Time     `jsonapi:"attr,live_starts_at,rfc3339,omitempty"`
	EndsAt        time.Time     `jsonapi:"attr,ends_at,rfc3339,omitempty"`
	LiveEndsAt    time.Time     `jsonapi:"attr,live_ends_at,rfc3339,omitempty"`
	LiveStartsAt  time.Time     `jsonapi:"attr,live_starts_at,rfc3339,omitempty"`
	TeamReminders []interface{} `jsonapi:"attr,team_reminders,rfc3339,omitempty"`
	TimeType      string        `jsonapi:"attr,time_type,omitempty"`
	UpdatedAt     time.Time     `jsonapi:"attr,updated_at,rfc3339,omitempty"`
	//relations
	AssignedTeams []Team `jsonapi:"relation,assigned_teams,omitempty"`
}
