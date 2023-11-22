package pco

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/db/models"
	"git.preston-baxter.com/Preston_PLB/capstone/webhook-service/vendors/pco/webhooks"
	"golang.org/x/oauth2"
)

var pcoMockAccount models.VendorAccount = models.VendorAccount{
	OauthCredentials: &models.OauthCredential{
		AccessToken:  "asdf;alskdfgha;dklrha;ldkfga;sldkf",
		ExpiresIn:    1234786012983,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		TokenType:    "bearer",
		RefreshToken: "asdfas;lkdjfas;dlkfj;asdlkj;aslf",
	},
	Name:             "pco",
}

func TestCreateSubscriptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Resp: %s", string(raw))
		w.Write([]byte(`{"data":[{"type":"Subscription","attributes":{"active":true,"name":"eventsandstuff","url":"https://thing.com/asdf/asdf/asdf"}},{"type":"Subscription","attributes":{"active":true,"name":"eventsandstuff","url":"https://thing.com/asdf/asdf/asdf"}}]}`))
	}))
	defer ts.Close()

	tokenSource := oauth2.StaticTokenSource(pcoMockAccount.Token())

	mockPco := config.VendorConfig{
		ClientId:      "as;dlkfja;slkdfj;aslkdfj;asdkl",
		ClientSecret:  "as;dlfkjas;ldkfja;slkdfj;alsdkfjas;dklj",
		Scopes:        []string{},
		AuthUri:       ts.URL,
		TokenUri:      ts.URL,
		RefreshEncode: "json",
		WebhookSecret: "as;dlfja;slkdja;slkdfj;alskdfj;alskdfa;slkdj",
	}

	pcoApi := NewClientWithOauthConfig(mockPco.OauthConfig(), tokenSource)
	if newUrl, err := url.Parse(ts.URL); err == nil {
		pcoApi.url = newUrl
	} else {
		t.Fatalf("%s", err)
	}

	mockSubscriptoins := []webhooks.Subscription{
		{
			Active:             true,
			Name:               "eventsandstuff",
			Url:                "https://thing.com/asdf/asdf/asdf",
		},
		{
			Active:             true,
			Name:               "eventsandstuff",
			Url:                "https://thing.com/asdf/asdf/asdf",
		},
	}

	_, err := pcoApi.CreateSubscriptions(mockSubscriptoins)
	if err != nil {
		t.Fatal(err)
	}


}
