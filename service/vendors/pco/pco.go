package pco

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const PCO_API_URL = "https://api.planningcenteronline.com"

type PcoApiClient struct {
	oauth *oauth2.Config
	token *oauth2.Token
	client *http.Client
	url *url.URL
}

func NewClient() *PcoApiClient {
	pco_url, err := url.Parse(PCO_API_URL)
	if err != nil {
		panic(err)
	}

	pco := &PcoApiClient{
		oauth: &oauth2.Config{},
		token: &oauth2.Token{},
		url: pco_url,
	}

	return pco
}

func NewClientWithOauthConfig(conf *oauth2.Config, token *oauth2.Token) *PcoApiClient {
	pco_url, err := url.Parse(PCO_API_URL)
	if err != nil {
		panic(err)
	}

	pco := &PcoApiClient{
		oauth: conf,
		token: token,
		url: pco_url,
	}

	return pco
}

func (api *PcoApiClient) getClient() *http.Client {
	if api.client == nil {
		api.client = api.oauth.Client(context.Background(), api.token)
	}

	return api.client
}

func (api *PcoApiClient) Url() *url.URL {
	return api.url
}

func (api *PcoApiClient) Do(req *http.Request) (*http.Response, error) {
	return api.getClient().Do(req)
}
