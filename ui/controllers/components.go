package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
)

func GetAddActionForm(c *gin.Context) {
	user := getUserFromContext(c)
	if user == nil {
		log.Warnf("Could not find user in context. Trying to redner Action form")
		badRequest(c, "No user available in context")
		return
	}

	accounts, err := mongo.FindVendorAccountByUser(user.Id)
	if err != nil {
		log.WithError(err).Errorf("Failed to find vendor accounts for: %s", user.Email)
		serverError(c, "No user available in context")
		return
	}

	renderTempl(c, templates.DashboardActionModal(accounts))
}
