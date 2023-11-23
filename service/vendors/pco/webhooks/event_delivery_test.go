package webhooks

import (
	"strings"
	"testing"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/services"
	"github.com/go-playground/assert/v2"
	"github.com/google/jsonapi"
)

func TestUnmarshallPayload(t *testing.T) {
	raw := `{"data":[{"id":"87f49852-1a2a-45cb-b4d2-0fa30eda0823","type":"EventDelivery","attributes":{"name":"services.v2.events.plan.created","attempt":1,"payload":"{\"data\":{\"type\":\"Plan\",\"id\":\"69259663\",\"attributes\":{\"can_view_order\":true,\"created_at\":\"2023-11-23T13:34:09Z\",\"dates\":\"No dates\",\"files_expire_at\":null,\"items_count\":0,\"last_time_at\":null,\"multi_day\":false,\"needed_positions_count\":0,\"other_time_count\":0,\"permissions\":\"Administrator\",\"plan_notes_count\":0,\"plan_people_count\":0,\"planning_center_url\":\"https://services.planningcenteronline.com/plans/69259663\",\"prefers_order_view\":false,\"public\":false,\"rehearsable\":true,\"rehearsal_time_count\":0,\"reminders_disabled\":false,\"series_title\":null,\"service_time_count\":0,\"short_dates\":\"No dates\",\"sort_date\":\"2023-11-23T13:34:09Z\",\"title\":null,\"total_length\":0,\"updated_at\":\"2023-11-23T13:34:09Z\"},\"relationships\":{\"service_type\":{\"data\":{\"type\":\"ServiceType\",\"id\":\"1429991\"}},\"next_plan\":{\"data\":null},\"previous_plan\":{\"data\":null},\"attachment_types\":{\"data\":[]},\"series\":{\"data\":null},\"created_by\":{\"data\":{\"type\":\"Person\",\"id\":\"136901110\"}},\"updated_by\":{\"data\":{\"type\":\"Person\",\"id\":\"136901110\"}},\"linked_publishing_episode\":{\"data\":null}},\"links\":{\"all_attachments\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/all_attachments\",\"attachments\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/attachments\",\"attendances\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/attendances\",\"contributors\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/contributors\",\"import_template\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/import_template\",\"item_reorder\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/item_reorder\",\"items\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/items\",\"live\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/live\",\"my_schedules\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/my_schedules\",\"needed_positions\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/needed_positions\",\"next_plan\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/next_plan\",\"notes\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/notes\",\"plan_times\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/plan_times\",\"previous_plan\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/previous_plan\",\"series\":null,\"signup_teams\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/signup_teams\",\"team_members\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663/team_members\",\"self\":\"https://api.planningcenteronline.com/services/v2/service_types/1429991/plans/69259663\",\"html\":\"https://services.planningcenteronline.com/plans/69259663\"}},\"included\":[],\"meta\":{\"can_include\":[\"contributors\",\"my_schedules\",\"plan_times\",\"series\"],\"parent\":{\"id\":\"1429991\",\"type\":\"ServiceType\"},\"event_time\":\"2023-11-23T13:34:09Z\"}}"},"relationships":{"organization":{"data":{"type":"Organization","id":"456240"}}}}]}`

	deliveries, err := jsonapi.UnmarshalManyPayload[EventDelivery](strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}

	plan := &services.Plan{}
	err = deliveries[0].UnmarshallPayload(plan)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, plan.Id, "69259663")
}
