package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"github.com/gin-gonic/gin"
)

const PCO_VALIDATE_HEADER = "X-PCO-Webhooks-Authenticity"

func ValidatePcoWebhook(c *gin.Context) {
	conf := config.Config()

	//get remote version from header
	remoteDigestStr := c.GetHeader(PCO_VALIDATE_HEADER)
	if remoteDigestStr == "" {
		log.Warnf("Request was sent with no %s header. Rejecting", PCO_VALIDATE_HEADER)
		c.AbortWithStatus(401)
		return
	}
	pcoSig := make([]byte, len(remoteDigestStr)/2)
	_, err := hex.Decode(pcoSig, []byte(remoteDigestStr))

	//clone request to harmlessly inspect the body
	bodyReader := c.Request.Clone(context.Background()).Body
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		log.WithError(err).Error("Failed to read body while validating PCO webhook")
		c.AbortWithError(501, err)
		return
	}

	//Get secret
	key := conf.Vendors[models.PCO_VENDOR_NAME].WebhookSecret

	//Get HMAC
	hmacSig := hmac.New(sha256.New, []byte(key))
	hmacSig.Write(body)

	if !hmac.Equal(hmacSig.Sum(nil), pcoSig) {
		log.Warn("")
		c.AbortWithStatus(401)
	}
}
