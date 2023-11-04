package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
)

type OauthCredential struct {
	AccessToken  string    `bson:"access_token,omitempty" json:"access_token,omitempty"`
	ExpiresIn    int       `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	ExpiresAt    time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	TokenType    string    `bson:"token_type,omitempty" json:"token_type,omitempty"`
	Scope        string    `bson:"scope,omitempty" json:"scope,omitempty"`
	RefreshToken string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
}

type OauthRefreshBody struct {
	ClientId     string `json:"cleint_id"`
	ClientSecret string `json:"cleint_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func (oc *OauthCredential) RefreshAccessToken(vendor string) error {
	conf := config.Config()
	vendorConfig := conf.Vendors[vendor]

	refresh_url, err := url.Parse(vendorConfig.TokenUri)
	if err != nil {
		return err
	}

	var body io.Reader
	switch vendorConfig.RefreshEncode {
	case "json":
		refreshBody := OauthRefreshBody{
			ClientId:     vendorConfig.ClientId,
			ClientSecret: vendorConfig.ClientSecret,
			GrantType:    "refresh_token",
			RefreshToken: oc.RefreshToken,
		}
		raw, err := json.Marshal(&refreshBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewReader(raw)
	case "url":
		q := refresh_url.Query()
		q.Add("client_id", vendorConfig.ClientId)
		q.Add("client_secret", vendorConfig.ClientSecret)
		q.Add("code", oc.RefreshToken)
		q.Add("grant_type", "refresh_token")

		body = strings.NewReader(q.Encode())
	default:
		panic(errors.New("Unkoown Encode Scheme"))
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", refresh_url.String(), body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	rawBody, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(rawBody, oc)
	if err != nil {
		return err
	}
	oc.ExpiresAt = time.Now().Add(time.Duration(oc.ExpiresIn)*time.Second - 10)

	return nil
}
