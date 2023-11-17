package controllers

import (
	"strings"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"github.com/gin-gonic/gin"
)

func AddActionFromForm(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		log.Warnf("Could not find user in context. Trying to redner Action form")
		badRequest(c, "No user available in context")
		return
	}
	//parse the form
	c.Request.ParseForm()
	var source []string
	var action []string

	//validate source
	if str := c.Request.FormValue("source"); str != "" {
		source = strings.Split(str, ".")
	} else {
		log.Warnf("Form request was partially or fully blank")
		badRequest(c, "Form request was partially or fully blank")
		return
	}

	//validate action
	if str := c.Request.FormValue("action"); str != "" {
		action = strings.Split(str, ".")
	} else {
		log.Warnf("Form request was partially or fully blank")
		badRequest(c, "Form request was partially or fully blank")
		return
	}

	//setup action

	//

	am := &models.ActionMapping{
		UserId: user.Id,
		SourceEvent: &models.Event{
			VendorName: source[0],
			Key:        source[1],
			Fields:     map[string]string{},
		},
		Action: &models.Action{
			VendorName: action[0],
			Type:       action[1],
			Fields:     map[string]string{},
		},
	}

	mongo.SaveModel(am)

	c.Redirect(302, "/dashboard")
}
