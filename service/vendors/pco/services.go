package pco

import (
	"fmt"
	"net/http"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/services"
	"github.com/google/jsonapi"
)

func (api *PcoApiClient) GetPlan(service_type_id, plan_id string) (*services.Plan, error){
	api.Url().Path = fmt.Sprintf("/services/v2/service_types/%s/plans/%s", service_type_id, plan_id)

	req, err := http.NewRequest(http.MethodGet, api.Url().String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("Failed to retrieve plan with status code: %d", resp.StatusCode)
	}

	plan := &services.Plan{}
	err = jsonapi.UnmarshalPayload(resp.Body, plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (api *PcoApiClient) GetPlanTimes(service_type_id, plan_id string) (*services.PlanTime, error) {
	api.Url().Path = fmt.Sprintf("/services/v2/service_types/%s/plans/%s/plan_times", service_type_id, plan_id)

	req, err := http.NewRequest(http.MethodGet, api.Url().String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("Failed to retrieve plan with status code: %d", resp.StatusCode)
	}

	planTime := &services.PlanTime{}
	err = jsonapi.UnmarshalPayload(resp.Body, planTime)
	if err != nil {
		return nil, err
	}

	return planTime, nil
}
