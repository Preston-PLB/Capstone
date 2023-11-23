package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
)

const PCO_VALIDATE_HEADER = "X-PCO-Webhooks-Authenticity"

func ValidatePcoWebhook(c *gin.Context) {
	//get remote version from header
	remoteDigestStr := c.GetHeader(PCO_VALIDATE_HEADER)
	if remoteDigestStr == "" {
		log.Warnf("Request was sent with no %s header. Rejecting", PCO_VALIDATE_HEADER)
		c.AbortWithStatus(401)
		return
	}
	pcoSig := make([]byte, len(remoteDigestStr)/2)
	_, err := hex.Decode(pcoSig, []byte(remoteDigestStr))
	if err != nil {
		log.WithError(err).Error("Failed to decode byte digest")
		_ = c.AbortWithError(501, err)
		return
	}

	//clone request body to harmlessly inspect the body
	bodyCopy := bytes.NewBuffer([]byte{})
	_, err = io.Copy(bodyCopy, c.Request.Body)
	if err != nil {
		log.WithError(err).Error("Failed to copy body while validating PCO webhook")
		_ = c.AbortWithError(501, err)
		return
	}
	body := bodyCopy.Bytes()
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	//Get secret
	key, err := getAuthSecret(c, body)
	if err != nil {
		log.WithError(err).Error("Failed to find auth secret for event. It may not be setup")
		_ = c.AbortWithError(501, err)
		return
	}

	//Get HMAC
	hmacSig := hmac.New(sha256.New, []byte(key))
	hmacSig.Write(body)

	if !hmac.Equal(hmacSig.Sum(nil), pcoSig) {
		log.Warn("")
		c.AbortWithStatus(401)
	}
}

func getAuthSecret(c *gin.Context, body []byte) (string, error) {
	userObjectId := userIdFromContext(c)
	log.Debug(string(body))

	//Pco is weird and sends a data array instead of an object. Yet there is only one event. Fun times
	event, err := jsonapi.UnmarshalManyPayload[webhooks.EventDelivery](bytes.NewBuffer(body))
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to unmarshall event delivery from PCO"), err)
	}

	if len(event) == 0 {
		return "", fmt.Errorf("There are no events in the delivery. Something is wrong")
	}

	webhook, err := mongo.FindPcoSubscriptionForUser(*userObjectId, event[0].Name)
	if err != nil {
		return "", errors.Join(fmt.Errorf("Failed to find pco subscription for user: %s and event: %s", userObjectId.Hex(), event[0].Name), err)
	}

	if webhook == nil {
		return "", fmt.Errorf("Could not find subscription for user: %s and name %s", userObjectId.Hex(), event[0].Name)
	}

	return webhook.Details.AuthenticitySecret, nil
}
