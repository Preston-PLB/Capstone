package controllers

import (
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/templates"
	"github.com/gin-gonic/gin"
)

func LandingPage(c *gin.Context) {
	if raw, exists := c.Get(USER_OBJ_KEY); exists {
		if user, ok := raw.(*models.User); ok {
			renderTempl(c, templates.LandingPage(user))
			return
		}
	}
	renderTempl(c, templates.LandingPage(nil))
}

func LoginPage(c *gin.Context) {
	renderTempl(c, templates.LoginPage(""))
}

func SignUpPage(c *gin.Context) {
	renderTempl(c, templates.SignupPage(""))
}

func DashboardPage(c *gin.Context) {
	if raw, exists := c.Get(USER_OBJ_KEY); exists {
		if user, ok := raw.(*models.User); ok {
			//get vendors
			vendors, err := mongo.FindVendorAccountByUser(user.MongoId())
			if err != nil {
				log.WithError(err).Error("Failed to get vendors for user: %s", user.Email)
				c.AbortWithStatus(502)
			}

			renderTempl(c, templates.DashboardPage(user, vendors))
			return
		} else {
			c.AbortWithStatus(502)
		}
	}
}
