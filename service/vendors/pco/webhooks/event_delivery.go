package webhooks

import (
	"strings"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/services"
	"github.com/google/jsonapi"
)

// Structure delivered to target when sending webhooks
type EventDelivery struct {
	//uuid of the EventDelivery
	ID string `jsonapi:"primary,EventDelivery"`
	//name of the event being sent. ex: services.v2.events.plan.updated
	//this coressponds to the scopes you set when configuring webhooks
	Name string `jsonapi:"attr,name"`
	//number of attemts taken to deliver the event
	Attempt int `jsonapi:"attr,attempt"`
	//JSON:API string of the event
	Payload string `jsonapi:"attr,payload"`
	//Owner Organization of the event
	Organization *services.Organization `jsonapi:"relation,organization"`
}

// Unmarshall payload of EventDelivery into the struct you think it is
func (event *EventDelivery) UnmarshallPayload(obj any) error {
	return jsonapi.UnmarshalPayload(strings.NewReader(event.Payload), obj)
}
