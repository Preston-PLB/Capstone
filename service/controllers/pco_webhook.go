package controllers

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	eventRegexKeys = map[string]string{"plan": `^services\.v2\.events\.plan\..*`}
	actionFuncMap = map[string]actionFunc{"youtube.livestream": ScheduleLiveStreamFromWebhook}
)

type actionFunc func(*gin.Context, *webhooks.EventDelivery) error

func ConsumePcoWebhook(c *gin.Context) {
	userId := c.Param("userid")

	if userId == "" {
		log.Warn("Webhook did not contain user id. Rejecting")
		c.AbortWithStatus(404)
		return
	}

	//get actions for user
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.WithError(err).Warn("User Id was malformed")
		c.AbortWithStatus(400)
		return
	}
	c.Set("user_bson_id", userObjectId)

	//read body and handle io in parallel because IO shenanigains
	wg := new(sync.WaitGroup)
	wg.Add(2)

	var actionMappings []models.ActionMapping
	var webhookBody *webhooks.EventDelivery
	errs := make([]error, 2)

	go func(wg *sync.WaitGroup) {
		actionMappings, errs[0] = mongo.FindActionMappingsByUser(userObjectId)
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		errs[1] = jsonapi.UnmarshalPayload(c.Request.Body, webhookBody)
		wg.Done()
	}(wg)

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		log.WithError(err).Errorf("Failed to do the IO parts")
		_ = c.AbortWithError(501, err)
		return
	}


	//perform actions
	//loop through all actions a user has
	for _, mapping := range actionMappings {
		//find the ones that are runable by this function
		if mapping.SourceEvent.VendorName == models.PCO_VENDOR_NAME && eventMatch(webhookBody.Name) {
			//generate lookup key for function
			actionKey := fmt.Sprintf("%s:%s", mapping.Action.VendorName, mapping.Action.Type)
			//if function exists run the function
			if action, ok := actionFuncMap[actionKey]; ok {
				err = action(c, webhookBody)
				//handle error
				if err != nil {
					log.WithError(err).Errorf("Failed to execute action: %s. From event source: %s:%s", actionKey, mapping.SourceEvent.VendorName, mapping.SourceEvent.Key)
					_ = c.AbortWithError(501, err)
				}
			}
		}
	}
}

func eventMatch(event string) bool {
	if regexString, ok := eventRegexKeys[event]; ok {
		reg := regexp.MustCompile(regexString)
		return reg.MatchString(event)
	} else {
		return false
	}
}

func youtubeServiceForUser(userId primitive.ObjectID) (*youtube.Service, error) {

	client, err := youtube.NewService(context.Background(), option.WithTokenSource(mongo.NewVendorTokenSource()))
	if err != nil {
		log.WithError(err).Error("Failed to init youtube service")
		panic(err)
	}
}

func ScheduleLiveStreamFromWebhook(c *gin.Context, body *webhooks.EventDelivery) error {

	return nil
}
