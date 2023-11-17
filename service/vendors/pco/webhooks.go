package pco

import (
	"fmt"
	"net/http"
	"reflect"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"github.com/google/jsonapi"
)


func (api *PcoApiClient) GetSubscriptions() ([]services.Subscription, error) {
	api.Url().Path = "/webhook/v2/subscriptions"

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

	raw, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(webhooks.Subscription{}))
	if err != nil {
		return nil, err
	}

	webhooks := make([]webhooks.Subscription, len(raw))
	for index, hook := range raw {
		var ok bool
		webhooks[index], ok = hook.(reflect.TypeOf(webhooks.Subscription))
		if !ok {
			return fmt.Errorf("Failed to extract webhook payload")
		}
	}



	return plan, nil
}
