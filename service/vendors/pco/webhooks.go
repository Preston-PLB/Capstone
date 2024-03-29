package pco

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"github.com/google/jsonapi"
)

// gets all current subscriptions
func (api *PcoApiClient) GetSubscriptions() ([]webhooks.Subscription, error) {
	api.Url().Path = "/webhooks/v2/subscriptions"

	req, err := http.NewRequest(http.MethodGet, api.Url().String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if raw, err := io.ReadAll(resp.Body); err == nil {
			return nil, fmt.Errorf("Failed to retrieve subscriptions with status code: %d. Error %s", resp.StatusCode, string(raw))
		} else {
			return nil, fmt.Errorf("Failed to retrieve subscriptions with status code: %d", resp.StatusCode)
		}
	}

	subscriptions, err := jsonapi.UnmarshalManyPayload[webhooks.Subscription](resp.Body)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Posts subscriptions to PCO api and returns a new list of subscriptions
func (api *PcoApiClient) CreateSubscriptions(subscriptions []webhooks.Subscription) ([]webhooks.Subscription, error) {
	api.Url().Path = "/webhooks/v2/subscriptions"

	body := bytes.NewBuffer([]byte{})
	err := jsonapi.MarshalPayload(body, subscriptions)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, api.Url().String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if raw, err := io.ReadAll(resp.Body); err == nil {
			return nil, fmt.Errorf("Failed to create subscriptions with status code: %d. Error %s", resp.StatusCode, string(raw))
		} else {
			return nil, fmt.Errorf("Failed to create subscriptions with status code: %d", resp.StatusCode)
		}
	}

	new_subscriptions, err := jsonapi.UnmarshalManyPayload[webhooks.Subscription](resp.Body)
	if err != nil {
		return nil, err
	}

	return new_subscriptions, nil
}

// Posts subcription to PCO api and updates the subscription at the pointer that was passed to the fuinction with the server response
func (api *PcoApiClient) CreateSubscription(subscription *webhooks.Subscription) error {
	api.Url().Path = "/webhooks/v2/subscriptions"

	body := bytes.NewBuffer([]byte{})
	err := jsonapi.MarshalPayload(body, subscription)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, api.Url().String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if raw, err := io.ReadAll(resp.Body); err == nil {
			return fmt.Errorf("Failed to create subscriptions with status code: %d. Error %s", resp.StatusCode, string(raw))
		} else {
			return fmt.Errorf("Failed to create subscription with status code: %d", resp.StatusCode)
		}
	}

	err = jsonapi.UnmarshalPayload(resp.Body, subscription)
	if err != nil {
		return err
	}

	return nil
}
