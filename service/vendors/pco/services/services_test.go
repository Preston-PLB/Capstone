package services_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/services"
	"github.com/go-playground/assert/v2"
	"github.com/google/jsonapi"
)

const valid_string = `{"data":{"type":"Plan","id":"69052110","attributes":{"can_view_order":true,"created_at":"2023-11-11T16:29:47Z","dates":"No dates","items_count":0,"multi_day":false,"needed_positions_count":0,"other_time_count":0,"permissions":"Administrator","plan_notes_count":0,"plan_people_count":0,"planning_center_url":"https://services.planningcenteronline.com/plans/69052110","prefers_order_view":true,"public":false,"rehearsable":true,"rehearsal_time_count":0,"reminders_disabled":false,"service_time_count":0,"short_dates":"No dates","sort_date":"2023-11-11T16:29:47Z","total_length":0,"updated_at":"2023-11-11T16:29:47Z"}}}`

func TestStructs(t *testing.T) {
	created_at, err := time.Parse(time.RFC3339, "2023-11-11T16:29:47Z")
	if err != nil {
		t.Fatal(err)
		return
	}

	sort_date, err := time.Parse(time.RFC3339, "2023-11-11T16:29:47Z")
	if err != nil {
		t.Fatal(err)
		return
	}

	updated_at, err := time.Parse(time.RFC3339, "2023-11-11T16:29:47Z")
	if err != nil {
		t.Fatal(err)
		return
	}

	plan := services.Plan{
		Id:                    "69052110",
		CanViewOrder:          true,
		CreatedAt:             created_at,
		Dates:                 "No dates",
		ItemsCount:            0,
		MultiDay:              false,
		NeededPositiionsCount: 0,
		OtherTimeCount:        0,
		Permissions:           "Administrator",
		PlanNotesCount:        0,
		PlanPeopleCount:       0,
		PlanningCenterUrl:     "https://services.planningcenteronline.com/plans/69052110",
		PerfersOrderView:      true,
		Public:                false,
		Rehearsable:           true,
		RehearsableTimeCount:  0,
		RemindersDisabled:     false,
		ServiceTimeCount:      0,
		ShortDates:            "No dates",
		Title:                 "",
		TotalLength:           0,
		SortDate:              sort_date,
		UpdatedAt:             updated_at,
	}

	valid_plan := &services.Plan{}
	test_plan := &services.Plan{}

	err = jsonapi.UnmarshalPayload(strings.NewReader(valid_string), valid_plan)
	if err != nil {
		t.Fatal(err)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	err = jsonapi.MarshalPayload(buf, &plan)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = jsonapi.UnmarshalPayload(buf, test_plan)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, test_plan, valid_plan)
}

func TestMarshalling(t *testing.T) {

}
