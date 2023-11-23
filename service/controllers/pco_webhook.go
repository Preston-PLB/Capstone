package controllers

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/services"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	yt_helpers "git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/youtube"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	eventRegexKeys = map[string]string{"plan": `^services\.v2\.events\.plan\..*`}
	actionFuncMap  = map[string]actionFunc{"youtube.livestream": ScheduleBroadcastFromWebhook}
)

type actionFunc func(*gin.Context, *webhooks.EventDelivery) error

func userIdFromContext(c *gin.Context) *primitive.ObjectID {
	if id, ok := c.Get("user_bson_id"); !ok {
		userId := c.Param("userid")

		if userId == "" {
			log.Warn("Webhook did not contain user id. Rejecting")
			c.AbortWithStatus(404)
			return nil
		}

		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.WithError(err).Warn("User Id was malformed")
			c.AbortWithStatus(400)
			return nil
		}
		c.Set("user_bson_id", userObjectId)
		return &userObjectId
	} else {
		if objId, ok := id.(primitive.ObjectID); ok {
			return &objId
		} else {
			return nil
		}
	}
}

func ConsumePcoWebhook(c *gin.Context) {
	userObjectId := userIdFromContext(c)

	//read body and handle io in parallel because IO shenanigains
	wg := new(sync.WaitGroup)
	wg.Add(2)

	//get actions for user
	var actionMappings []models.ActionMapping
	var webhookBody *webhooks.EventDelivery
	errs := make([]error, 2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		actionMappings, errs[0] = mongo.FindActionMappingsByUser(*userObjectId)
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var payload []webhooks.EventDelivery
		payload, errs[1] = jsonapi.UnmarshalManyPayload[webhooks.EventDelivery](c.Request.Body)
		webhookBody = &payload[0]
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
		if mapping.SourceEvent.VendorName == models.PCO_VENDOR_NAME && eventMatch(mapping.SourceEvent.Key, webhookBody.Name) {
			//generate lookup key for function
			actionKey := fmt.Sprintf("%s.%s", mapping.Action.VendorName, mapping.Action.Type)
			//if function exists run the function
			if action, ok := actionFuncMap[actionKey]; ok {
				err := action(c, webhookBody)
				//handle error
				if err != nil {
					log.WithError(err).Errorf("Failed to execute action: %s. From event source: %s:%s", actionKey, mapping.SourceEvent.VendorName, mapping.SourceEvent.Key)
					_ = c.AbortWithError(501, err)
				} else {
					log.Infof("Succesfully proccessed: %s for %s", webhookBody.Name, userObjectId.Hex())
					c.Status(200)
				}
				return
			}
		}
	}
	log.Warnf("No errors, but also no work...")
	c.Status(200)
}

func eventMatch(key, event string) bool {
	if regexString, ok := eventRegexKeys[key]; ok {
		reg := regexp.MustCompile(regexString) //TODO: Make this regex cache-able
		return reg.MatchString(event)
	} else {
		return false
	}
}

func pcoServiceForUser(userId primitive.ObjectID) (*pco.PcoApiClient, error) {
	//add youtube client to map if its not there
	if client, ok := pcoClientMap[userId]; !ok {
		pcoAccount, err := mongo.FindVendorAccountByUser(userId, models.PCO_VENDOR_NAME)
		if err != nil {
			return nil, err
		}

		//Build our fancy token source
		tokenSource := oauth2.ReuseTokenSource(pcoAccount.Token(), mongo.NewVendorTokenSource(pcoAccount))

		//init service
		conf := config.Config()
		client := pco.NewClientWithOauthConfig(conf.Vendors[models.PCO_VENDOR_NAME].OauthConfig(), tokenSource)

		//add user to map
		pcoClientMap[userId] = client

		return client, nil
	} else {
		return client, nil
	}
}

func youtubeServiceForUser(userId primitive.ObjectID) (*youtube.Service, error) {
	//add youtube client to map if its not there
	if client, ok := ytClientMap[userId]; !ok {
		ytAccount, err := mongo.FindVendorAccountByUser(userId, models.YOUTUBE_VENDOR_NAME)
		if err != nil {
			return nil, err
		}

		//Build our fancy token source
		tokenSource := oauth2.ReuseTokenSource(ytAccount.Token(), mongo.NewVendorTokenSource(ytAccount))

		//init service
		client, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource))
		if err != nil {
			log.WithError(err).Error("Failed to init youtube service")
			return nil, err
		}

		//add user to map
		ytClientMap[userId] = client

		return client, nil
	} else {
		return client, nil
	}
}

func ScheduleBroadcastFromWebhook(c *gin.Context, body *webhooks.EventDelivery) error {
	//get uid from context. Lots of sanitizing just incase
	uid := userIdFromContext(c)

	//Check if this is a redilivery.

	//Load ytClient for user. It is fetched from cache or created
	ytClient, err := youtubeServiceForUser(*uid)
	if err != nil {
		log.WithError(err).Error("Failed to initialize youtube client")
		return err
	}

	//Load pcoClient for user. It is fetched from cache or created
	pcoClient, err := pcoServiceForUser(*uid)
	if err != nil {
		log.WithError(err).Error("Failed to initialize youtube client")
		return err
	}

	//deserialize the payload
	payload := &services.Plan{}
	err = body.UnmarshallPayload(payload)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshall body")
		return err
	}

	//Save audit point
	eventRecievedAudit := &models.EventRecieved{
		UserId:     *uid,
		VendorName: models.PCO_VENDOR_NAME,
		VendorId:   body.ID,
		Type:       body.Name,
	}

	if err := mongo.SaveModel(eventRecievedAudit); err != nil {
		log.WithError(err).WithField("EventRecieved", eventRecievedAudit).Error("Failed to save audit event. Logging here and resuming")
	}

	//create the broadcast
	//TODO: handle update
	//TODO: handle delete
	var broadcast *youtube.LiveBroadcast
	switch body.Name {
	case "services.v2.events.plan.created":
		broadcast, err = scheduleNewBroadcastFromWebhook(c, payload, ytClient, pcoClient)
		if err != nil {
			log.WithError(err).Error("Failed to schedule broadcast from created event")
			return err
		}
	case "services.v2.events.plan.updated":
		log.Warn("services.v2.events.plan.updated event not implemented")
		return nil
	case "services.v2.events.plan.destroyed":
		log.Warn("services.v2.events.plan.destroyed event not implemented")
		return nil
	default:
		return fmt.Errorf("Unkown event error: %s", body.Name)
	}

	//build audit trail after action was taken
	broadcastModel := &models.YoutubeBroadcast{
		UserId:  *uid,
		Details: broadcast,
	}

	actionTaken := &models.ActionTaken{
		UserId:          *uid,
		TriggeringEvent: eventRecievedAudit.MongoId(),
		Result:          []primitive.ObjectID{broadcastModel.MongoId()},
		VendorName:      models.YOUTUBE_VENDOR_NAME,
	}

	//save audit trail
	err = mongo.SaveModels(broadcastModel, actionTaken)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshall body")
		return err
	}

	return nil
}

func scheduleNewBroadcastFromWebhook(c *gin.Context, plan *services.Plan, ytClient *youtube.Service, pcoClient *pco.PcoApiClient) (*youtube.LiveBroadcast, error) {
	times, err := pcoClient.GetPlanTimes(plan.ServiceType.Id, plan.Id)
	if err != nil {
		return nil, err
	}

	startTime := times[0].StartsAt
	// endTime := times[len(times) - 1].EndsAt TODO: this will be used later
	var title string
	if plan.Title == "" {
		title = "Live Stream Scheduled By Capstone"
	} else {
		title = plan.Title
	}

	return yt_helpers.InsertBroadcast(ytClient, title, startTime, yt_helpers.STATUS_PRIVATE)
}
