package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"github.com/gin-gonic/gin"
)

const REDIRECT_URI = "https://capstone.preston-baxter.com:8080/vendor/youtube/callback"

func InitiateYoutubeOuath(c *gin.Context) {
	conf := config.Config()

	init_url, err := url.Parse(conf.YoutubeConfig.AuthUri)
	if err != nil {
		//we should not get here
		panic(err)
	}

	q := init_url.Query()
	//https://developers.google.com/youtube/v3/guides/auth/server-side-web-apps#httprest_1
	q.Add("client_id", conf.YoutubeConfig.ClientId)
	q.Add("redirect_uri", REDIRECT_URI)
	q.Add("response_type", "code")
	q.Add("scope", conf.YoutubeConfig.Scope())
	q.Add("access_type", "offline")
	//used to prevent CSRF
	q.Add("state", getAuthHash(c))
	init_url.RawQuery = q.Encode()

	c.Redirect(302, init_url.String())
}

func ReceiveYoutubeOauth(c *gin.Context) {
	conf := config.Config()
	user := getUserFromContext(c)

	if user == nil {
		log.Error("Unable to find user in context")
		c.AbortWithStatus(502)
	}

	code := c.Query("code")
	respHash := c.Query("state")
	//validate returned code
	if code == "" {
		log.Error("Youtube OAuth response did not contain a code. Possible CSRF")
		c.AbortWithStatus(502)
		return
	}

	//validate state
	if respHash == "" || respHash != getAuthHash(c) {
		log.Error("Youtube OAuth response did not contain the correct hash. Possible CSRF")
		c.AbortWithStatus(502)
		return
	}

	client := http.Client{}

	token_url, err := url.Parse(conf.YoutubeConfig.TokenUri)
	if err != nil {
		//we should not get here
		panic(err)
	}

	//Make request to google for credentials
	q := token_url.Query()

	q.Add("code", code)
	q.Add("client_id", conf.YoutubeConfig.ClientId)
	q.Add("client_secret", conf.YoutubeConfig.ClientSecret)
	q.Add("redirect_uri", REDIRECT_URI)
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
		log.Errorf("Response failed with status code: %d. Error: %s", resp.StatusCode ,string(rawBody))
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
		Name:             "youtube",
	}

	err = mongo.SaveModel(vendor)
	if err != nil {
		log.WithError(err).Errorf("Failed to save credentials for user: %s", user.Email)
		c.AbortWithStatus(502)
		return
	}

	c.Redirect(302, "/dashboard")
}
