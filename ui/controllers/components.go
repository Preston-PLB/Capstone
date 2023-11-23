package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetAddActionForm(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		log.Warnf("Could not find user in context. Trying to redner Action form")
		badRequest(c, "No user available in context")
		return
	}

	accounts, err := mongo.FindAllVendorAccountsByUser(user.Id)
	if err != nil {
		log.WithError(err).Errorf("Failed to find vendor accounts for: %s", user.Email)
		serverError(c, "No user available in context")
		return
	}

	renderTempl(c, templates.DashboardActionModal(accounts))
}

type DashboardMetric struct {
	Title          string
	PrimaryValue   string
	SecondaryValue string
	Subtitle       string
}

type dashboardMetricFunc func(c *gin.Context) *DashboardMetric

var metricFuncMap = map[string]dashboardMetricFunc{"default": defaultMetricFunction, "events_received": eventsRecievedMetricFunction, "streams_scheduled": streamsScheduledMetricFunction}

func GetMetricCard(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		log.Warnf("Could not find user in context. Trying to redner Action form")
		badRequest(c, "No user available in context")
		return
	}

	if metric, ok := c.GetQuery("metric"); ok {
		if metricFunc, mok := metricFuncMap[metric]; mok {
			renderDashboardMetric(c, &metricFunc)
			return
		}
	}

	//send default metric function
	log.Warn("Failed to find metricfunc")
	defaultFunc := metricFuncMap["default"]
	renderDashboardMetric(c, &defaultFunc)
}

func defaultMetricFunction(c *gin.Context) *DashboardMetric {
	return &DashboardMetric{
		Title:          "Err",
		PrimaryValue:   "0.00",
		SecondaryValue: "0.00",
		Subtitle:       "something went wrong",
	}
}

func renderDashboardMetric(c *gin.Context, metricFunc *dashboardMetricFunc) {
	metric := (*metricFunc)(c)
	renderTempl(c, templates.DashboardCard(metric.Title, metric.PrimaryValue, metric.SecondaryValue, metric.Subtitle))
}

func eventsRecievedMetricFunction(c *gin.Context) *DashboardMetric {
	user := getUserFromContext(c)

	events, err := mongo.AggregateVendorEventReport(user.Id)
	if err != nil {
		log.WithError(err).Errorf("Failed to find events for user: %s", user.Id.Hex())
		return defaultMetricFunction(c)
	}

	totalEvents := 0
	biggestVendor := 0
	for index, event := range events {
		totalEvents += event.Count
		if events[biggestVendor].Count < event.Count {
			biggestVendor = index
		}

	}

	p := message.NewPrinter(language.English)
	return &DashboardMetric{
		Title:          "Events Recieved",
		PrimaryValue:   p.Sprintf("%d", totalEvents),
		SecondaryValue: p.Sprintf("Most events came from: %s", events[biggestVendor].Name),
		Subtitle:       "thats a lot of events",
	}
}

func streamsScheduledMetricFunction(c *gin.Context) *DashboardMetric {
	return defaultMetricFunction(c)
}
