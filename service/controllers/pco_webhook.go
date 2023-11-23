package controllers

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db"
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

	//Error definintions
	errorOkMap                = map[error]bool{NotSchedulableTime: true, UnknownEventErr: true, AlreadyScheduledBroadcast: true, NoBroadcastToDelete: true}
	NotSchedulableTime        = errors.New("This time is not schedulable")
	UnknownEventErr           = errors.New("Event sent is unkown")
	AlreadyScheduledBroadcast = errors.New("This broadcast has already been scheduled")
	NoBroadcastToDelete       = errors.New("No Broadcasts to destroy")
)

const (
	CREATED_BROADCAST = "Created Broadcast"
	UPDATED_BROADCAST = "Updated Broadcast"
	DELETED_BROADCAST = "Deleted Broadcast"
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
					//if err is in the ok map, return 200
					if pass, ok := errorOkMap[err]; ok && pass {
						log.Warnf("Continueing after error: %s. From action: %s. From event source: %s:%s", err, actionKey, mapping.SourceEvent.VendorName, mapping.SourceEvent.Key)
						c.Status(200)
					} else {
						log.WithError(err).Errorf("Failed to execute action: %s. From event source: %s:%s", actionKey, mapping.SourceEvent.VendorName, mapping.SourceEvent.Key)
						_ = c.AbortWithError(501, err)
					}
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

//TODO: Revisit the structure of this function
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

	//Check to see if we have scheduled a broadcast befre
	broadcasts, err := mongo.FindAllBroadcastsByCorrelationId(*uid, payload.Id)
	if err != nil {
		return errors.Join(fmt.Errorf("Failed to find broadcasts for user: %s and CorrelationId: %s", uid.Hex(), payload.Id), err)
	}

	var result string
	if len(broadcasts) > 0 {
		//What do we do when we have already scheduled the broadcast
		switch body.Name {
		//If we get plan created event for this, return already scheduled error
		case "services.v2.events.plan.created":
			return AlreadyScheduledBroadcast
		//update the broadcast
		case "services.v2.events.plan.updated":
			//TODO: Update Broadcast
			err := updateBroadcastFromWebhook(c, broadcasts, payload, ytClient, pcoClient)
			if err != nil {
				log.WithError(err).Error("Failed to update broadcast from updated event")
				return err
			}
			result = UPDATED_BROADCAST
		//delete the broadcast
		case "services.v2.events.plan.destroyed":
			//TODO: Delete broadcast
			err := deleteBroadcastFromWebhook(c, broadcasts, payload, ytClient, pcoClient)
			if err != nil {
				log.WithError(err).Error("Failed to delete broadcast from updated event")
				return err
			}
			result = DELETED_BROADCAST
		default:
			return UnknownEventErr
		}
		actionTaken := &models.ActionTaken{
			UserId:          *uid,
			TriggeringEvent: eventRecievedAudit.MongoId(),
			Result:          result,
			VendorName:      models.YOUTUBE_VENDOR_NAME,
		}

		//save audit trail
		err = mongo.SaveModels( actionTaken)
		if err != nil {
			log.WithError(err).Error("Failed to save broadcastModel and actionTaken")
			return err
		}
	} else {
		//No broadcast is scheduled
		//create the broadcast
		var broadcast *youtube.LiveBroadcast
		switch body.Name {
		case "services.v2.events.plan.created":
			broadcast, err = scheduleNewBroadcastFromWebhook(c, payload, ytClient, pcoClient)
			if err != nil {
				log.WithError(err).Error("Failed to schedule broadcast from created event")
				return err
			}
			result = CREATED_BROADCAST
		case "services.v2.events.plan.updated":
			broadcast, err = scheduleNewBroadcastFromWebhook(c, payload, ytClient, pcoClient)
			if err != nil {
				log.WithError(err).Error("Failed to schedule broadcast from updated event")
				return err
			}
			result = CREATED_BROADCAST
		case "services.v2.events.plan.destroyed":
			return NoBroadcastToDelete
		default:
			return fmt.Errorf("Unkown event error: %s", body.Name)
		}

		//build audit trail after action was taken
		broadcastModel := &models.YoutubeBroadcast{
			UserId:        *uid,
			CorrelationId: payload.Id,
			Details:       broadcast,
		}

		actionTaken := &models.ActionTaken{
			UserId:          *uid,
			TriggeringEvent: eventRecievedAudit.MongoId(),
			Result:          result,
			VendorName:      models.YOUTUBE_VENDOR_NAME,
		}

		//save audit trail
		err = mongo.SaveModels(broadcastModel, actionTaken)
		if err != nil {
			log.WithError(err).Error("Failed to save broadcastModel and actionTaken")
			return err
		}
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

	//if starttime is before now, skip with a passable error
	if startTime.Before(time.Now()) {
		return nil, NotSchedulableTime
	}

	var title string
	if plan.Title == "" {
		title = "Live Stream Scheduled By Capstone"
	} else {
		title = plan.Title
	}

	log.Debugf("Trying to schedule time at: %s", startTime.Format(yt_helpers.ISO_8601))
	return yt_helpers.InsertBroadcast(ytClient, title, startTime, yt_helpers.STATUS_PRIVATE)
}

func updateBroadcastFromWebhook(c *gin.Context, broadcasts []models.YoutubeBroadcast, plan *services.Plan, ytClient *youtube.Service, pcoClient *pco.PcoApiClient) error {
	times, err := pcoClient.GetPlanTimes(plan.ServiceType.Id, plan.Id)
	if err != nil {
		return err
	}
	startTime := times[0].StartsAt
	// endTime := times[len(times) - 1].EndsAt TODO: this will be used later

	//if starttime is before now, skip with a passable error
	if startTime.Before(time.Now()) {
		return NotSchedulableTime
	}

	var title string
	if plan.Title == "" {
		title = "Live Stream Scheduled By Capstone"
	} else {
		title = plan.Title
	}

	//create list of errors to process all of the broadcasts and then error
	errs := make([]error, 0, len(broadcasts))
	bcs := make([]*models.YoutubeBroadcast, 0, len(broadcasts))
	for index, broadcast := range broadcasts {
		liveBroadcast, err := yt_helpers.UpdateBroadcast(ytClient, broadcast.Details.Id, title, startTime, yt_helpers.STATUS_PRIVATE)
		if err != nil {
			errs = append(errs, err)
		} else {
			broadcasts[index].Details = liveBroadcast
			bcs = append(bcs, &broadcasts[index])
		}
	}

	if err := errors.Join(errs...); err != nil {
		return err
	}

	return db.SaveModelSlice(mongo, bcs...)
}

func deleteBroadcastFromWebhook(c *gin.Context, broadcasts []models.YoutubeBroadcast, plan *services.Plan, ytClient *youtube.Service, pcoClient *pco.PcoApiClient) error {
	errs := make([]error, 0, len(broadcasts))
	bcs := make([]*models.YoutubeBroadcast, 0, len(broadcasts))
	for index, broadcast := range broadcasts {
		err := yt_helpers.DeleteBroadcast(ytClient, broadcast.Details.Id)
		if err != nil {
			errs = append(errs, err)
		} else {
			bcs = append(bcs, &broadcasts[index])
		}
	}

	if err := errors.Join(errs...); err != nil {
		return err
	}

	return db.DeleteModelSlice(mongo, bcs...)
}
