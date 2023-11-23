package controllers

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPlanEventMatch(t *testing.T) {
	events := []string{"services.v2.events.plan.updated", "services.v2.events.plan.destroyed", "services.v2.events.plan.created"}

	for _, event := range events {
		assert.Equal(t, eventMatch("plan", event), true)
	}
}
