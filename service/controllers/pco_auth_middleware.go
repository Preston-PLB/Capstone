package controllers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

	//clone request to harmlessly inspect the body
	bodyReader := c.Request.Clone(context.Background()).Body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		log.WithError(err).Error("Failed to read body while validating PCO webhook")
		_ = c.AbortWithError(501, err)
		return
	}

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

	event := &webhooks.EventDelivery{}
	err := jsonapi.UnmarshalPayload(bytes.NewBuffer(body), event)
	if err != nil {
		return "", err
	}

	webhook, err := mongo.FindPcoSubscriptionForUser(*userObjectId, event.Name)
	if err != nil {
		return "", err
	}

	return webhook.Details.AuthenticitySecret, nil
}
