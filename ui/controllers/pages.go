package controllers

import (
	"errors"
	"sync"

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
	user := getUserFromContext(c)

	if user == nil {
		log.Error("No user found in context")
		serverError(c, "No user found in context")
		return
	}

	//Split database fetching into go routines
	var vendors []models.VendorAccount
	var actions []models.ActionMapping
	//TODO: find a generic way to do this.
	errs := make([]error, 2)

	//Use waitgroup to syncronize
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go func(wg *sync.WaitGroup) {
		vendors, errs[0] = mongo.FindAllVendorAccountsByUser(user.MongoId())
		wg.Done()
	}(waitGroup)

	go func(wg *sync.WaitGroup) {
		actions, errs[1] = mongo.FindActionMappingsByUser(user.MongoId())
		wg.Done()
	}(waitGroup)

	//after this line we are in sync
	waitGroup.Wait()

	//handle errors
	for _, err := range errs {
		if err != nil {
			log.WithError(errors.Join(errs...)).Error("Failed to do database lookup when retrieving dashbDashboardPage")
			serverError(c, "Failed to do database lookup when retrieving dashbDashboardPage")
			return
		}
	}

	renderTempl(c, templates.DashboardPage(user, vendors, actions))
}

func EventsPage(c *gin.Context) {
	user := getUserFromContext(c)

	if user == nil {
		log.Error("No user found in context")
		serverError(c, "No user found in context")
		return
	}

	renderTempl(c, templates.EventsPage(user))
}
