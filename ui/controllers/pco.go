package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"github.com/gin-gonic/gin"
)

const PCO_REDIRECT_URI = "https://%s/vendor/pco/callback"

func InitiatePCOOuath(c *gin.Context) {
	conf := config.Config()
	vendorConfig := conf.Vendors[models.PCO_VENDOR_NAME]

	init_url, err := url.Parse(vendorConfig.AuthUri)
	if err != nil {
		//we should not get here
		panic(err)
	}

	q := init_url.Query()
	q.Add("client_id", vendorConfig.ClientId)
	q.Add("redirect_uri", fmt.Sprintf(PCO_REDIRECT_URI, conf.AppSettings.FrontendServiceUrl))
	q.Add("response_type", "code")
	q.Add("scope", vendorConfig.Scope())
	init_url.RawQuery = q.Encode()

	c.Redirect(302, init_url.String())
}

func RecievePCOOuath(c *gin.Context) {
	conf := config.Config()
	vendorConfig := conf.Vendors[models.PCO_VENDOR_NAME]
	user := getUserFromContext(c)

	if user == nil {
		log.Error("Unable to find user in context")
		c.AbortWithStatus(502)
	}

	code := c.Query("code")
	//validate returned code
	if code == "" {
		log.Error("Youtube OAuth response did not contain a code. Possible CSRF")
		c.AbortWithStatus(502)
		return
	}

	client := http.Client{}

	token_url, err := url.Parse(vendorConfig.TokenUri)
	if err != nil {
		//we should not get here
		panic(err)
	}

	//Make request to google for credentials
	q := token_url.Query()

	q.Add("code", code)
	q.Add("client_id", vendorConfig.ClientId)
	q.Add("client_secret", vendorConfig.ClientSecret)
	q.Add("redirect_uri", fmt.Sprintf(PCO_REDIRECT_URI, conf.AppSettings.FrontendServiceUrl))
	q.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", token_url.String(), strings.NewReader(q.Encode()))
	if err != nil {
		log.WithError(err).Errorf("Failed to generate request with the following url: '%s'", token_url.String())
		c.AbortWithStatus(502)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Errorf("Failed to make request to the following url: '%s'", token_url.String())
		c.AbortWithStatus(502)
		return
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("Failed to read body from the following url: '%s'", token_url.String())
		c.AbortWithStatus(502)
		return
	}

	if resp.StatusCode != 200 {
		log.Errorf("Response failed with status code: %d. Error: %s", resp.StatusCode, string(rawBody))
		c.AbortWithStatus(502)
		return
	}

	oauthResp := &models.OauthCredential{}
	err = json.Unmarshal(rawBody, oauthResp)
	if err != nil {
		log.WithError(err).Errorf("Failed to Unmarshal response from the following url: '%s'", token_url.String())
		c.AbortWithStatus(502)
	}
	log.Infof("oauthResp: %v", *oauthResp)
	//Set expires at time but shave some time off to refresh token before expire date
	oauthResp.ExpiresAt = time.Now().Add(time.Duration(oauthResp.ExpiresIn)*time.Second - 10)

	//store credentials
	vendor := &models.VendorAccount{
		UserId:           user.Id,
		OauthCredentials: oauthResp,
		Name:             models.PCO_VENDOR_NAME,
	}

	err = mongo.SaveModel(vendor)
	if err != nil {
		log.WithError(err).Errorf("Failed to save credentials for user: %s", user.Email)
		c.AbortWithStatus(502)
		return
	}

	c.Redirect(302, "/dashboard")

}
