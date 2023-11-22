package controllers

import (
	"errors"
	"fmt"
	"strings"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type actionFunc func(user *models.User) error

var (
	actionFuncs      map[string]actionFunc            = map[string]actionFunc{"pco.plan": setupPcoSubscriptions}
	webhooksTemplate map[string]webhooks.Subscription = map[string]webhooks.Subscription{
		"services.v2.events.plan.created": {
			Active: true,
			Name:   "services.v2.events.plan.created",
			Url:    "https://%s/pco/%s",
		},
		"services.v2.events.plan.updated": {
			Active: true,
			Name:   "services.v2.events.plan.updated",
			Url:    "https://%s/pco/%s",
		},
		"services.v2.events.plan.deleted": {
			Active: true,
			Name:   "services.v2.events.plan.deleted",
			Url:    "https://%s/pco/%s",
		},
	}
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

	//setup action listener
	if afunc, ok := actionFuncs[strings.Join(source, ".")]; ok {
		err := afunc(user)
		if err != nil {
			log.WithError(err).Error("Failed to setup actions")
			serverError(c, "Failed to setup actions")
			return
		}
	}

	//Build mappings
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

	err := mongo.SaveModel(am)
	if err != nil {
		log.WithError(err).Error("Failed to setup actions")
		serverError(c, "Failed to setup actions")
		return
	}

	c.Redirect(302, "/dashboard")
}

func setupPcoSubscriptions(user *models.User) error {
	// Get PCO vendor account
	conf := config.Config()
	pcoAccount, err := mongo.FindVendorAccountByUser(user.Id, models.PCO_VENDOR_NAME)
	if err != nil {
		return err
	}

	//build pco api
	tokenSource := oauth2.ReuseTokenSource(pcoAccount.Token(), mongo.NewVendorTokenSource(pcoAccount))
	pcoApi := pco.NewClientWithOauthConfig(conf.Vendors[models.PCO_VENDOR_NAME].OauthConfig(), tokenSource)

	//Check if subscriptions already exist
	webhookMap := make(map[string]webhooks.Subscription)
	subscriptions, err := pcoApi.GetSubscriptions()
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to find subscriptions for user: %s", user.Id), err)
	}
	//Loop through found subscriptions
	for _, sub := range subscriptions {
		//if subsciption is in the templates look to add it to our map
		if templ, ok := webhooksTemplate[sub.Name]; ok {
			//if the subscription is for our url add it to our map
			url := fmt.Sprintf(templ.Url, conf.AppSettings.WebhookServiceUrl, user.Id.Hex())
			if url == sub.Url {
				webhookMap[sub.Name] = sub
			}
		}
	}

	builtHooks := make([]webhooks.Subscription, 0, len(webhooksTemplate))
	//Build subscriptions
	for _, templ := range webhooksTemplate {
		if _, ok := webhookMap[templ.Name]; !ok {
			builtHooks = append(builtHooks, webhooks.Subscription{
				Active: false,
				Name:   templ.Name,
				Url:    fmt.Sprintf(templ.Url, conf.AppSettings.WebhookServiceUrl, user.Id.Hex()),
			})
		}
	}

	//Post Subscriptions
	subscriptions, err = pcoApi.CreateSubscriptions(builtHooks)
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to create subscriptions for user: %s", user.Id), err)
	}

	//Save Subscriptions
	err = mongo.SaveSubscriptionsForUser(user.Id, subscriptions...)
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to save subscriptions for user: %s", user.Id), err)
	}

	return nil
}
