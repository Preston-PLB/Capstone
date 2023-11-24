package controllers

import (
	"strings"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
)

var (
	eventsTableMap = map[string]eventsTableFunc{"default": defaultTableData, "events_for_user": eventsForUserTableData, "actions_for_user": actionsForUserTableData}
)

type eventsTableFunc func(c *gin.Context) templates.TableData

func GetTableComponent(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		log.Warnf("Could not find user in context. Trying to redner Action form")
		badRequest(c, "No user available in context")
		return
	}

	if table, ok := c.GetQuery("table_name"); ok {
		if tableFunc, mok := eventsTableMap[table]; mok {
			renderEventTable(c, table, &tableFunc)
			return
		}
	}

	//send default metric function
	log.Warn("Failed to find eventsTableFunc")
	defaultFunc := eventsTableMap["default"]
	renderEventTable(c, "default", &defaultFunc)
}

func defaultTableData(c *gin.Context) templates.TableData {
	return [][]string{{"id", "col 1", "col 2"}, {"row data", "item data", "stuff"}}
}

func renderEventTable(c *gin.Context, title string, tableFunc *eventsTableFunc) {
	tableData := (*tableFunc)(c)
	renderTempl(c, templates.EventTableData(tableData, title))
}

func eventsForUserTableData(c *gin.Context) templates.TableData {
	//User can't be nil because we check before we get here
	user := getUserFromContext(c)

	events, err := mongo.FindEventsRecievedByUserId(user.Id)
	if err != nil {
		log.WithError(err).Errorf("Failed to find events for user: %s to load table data", user.Id.Hex())
		return defaultTableData(c)
	}

	//check for filter
	filter, filter_exists := c.GetQuery("filter")

	table := make([][]string, len(events)+1)
	index := 1
	for _, event := range events {
		arr := []string{event.CreatedAt.Format(time.Stamp), strings.ToUpper(event.VendorName), event.VendorId, event.Type}

		if filter_exists {
			//if the filter exists loop through the row. Check if anything meets the filter
			pass := false
			for _, item := range arr {
				//if we already have a match short circuit. If we don't we can potentially flip from false -> true
				pass = pass || strings.Contains(item, filter)
			}
			//If we did not find a matching item continue
			if !pass {
				continue
			}
		}
		//We either had no filter or passed the filter check. Add to the pool
		table[index] = arr
		index += 1
		table[0] = []string{"Timestamp", "Vendor", "Id", "Event Type"}
	}

	return table[0:index]
}

func actionsForUserTableData(c *gin.Context) templates.TableData {
	//User can't be nil because we check before we get here
	user := getUserFromContext(c)

	actions, err := mongo.FindActionTakenByUserId(user.Id)
	if err != nil {
		log.WithError(err).Errorf("Failed to find actions for user: %s to load table data", user.Id.Hex())
		return defaultTableData(c)
	}


	//check for filter
	filter, filter_exists := c.GetQuery("filter")
	index := 1
	table := make([][]string, len(actions)+1)
	for _, action := range actions {
		arr := []string{action.CreatedAt.Format(time.RFC1123), action.VendorName, action.CorrelationId, action.Result}
		if filter_exists {
			//if the filter exists loop through the row. Check if anything meets the filter
			pass := false
			for _, item := range arr {
				//if we already have a match short circuit. If we don't we can potentially flip from false -> true
				pass = pass || strings.Contains(item, filter)
			}
			//If we did not find a matching item continue
			if !pass {
				continue
			}
		}
		table[index] = []string{action.CreatedAt.Format(time.RFC1123), action.VendorName, action.CorrelationId, action.Result}
		index += 1
	}
	table[0] = []string{"Timestamp", "Vendor", "Id", "Result"}

	return table[0:index]
}
