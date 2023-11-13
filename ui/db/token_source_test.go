package db

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"git.preston-baxter.com/Preston_PLB/capstone/frontend-service/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

func newConf(url string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		RedirectURL:  "REDIRECT_URL",
		Scopes:       []string{"scope1", "scope2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  url + "/auth",
			TokenURL: url + "/token",
		},
	}
}

func TestRefreshToken(t *testing.T) {
	config.Init()
	conf := config.Config()

	client1, err := NewClient(conf.Mongo.Uri)
	if err != nil {
		t.Fatal(err)
	}

	client2, err := NewClient(conf.Mongo.Uri)
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"ACCESS_TOKEN",  "scope": "user", "token_type": "bearer", "refresh_token": "NEW_REFRESH_TOKEN"}`))
		count += 1
	}))
	defer ts.Close()

	conf.Vendors["test"].AuthUri = fmt.Sprintf("%s/auth", ts.URL)
	conf.Vendors["test"].TokenUri = fmt.Sprintf("%s/auth", ts.URL)
	//tkr := conf.TokenSource(context.Background(), &oauth2.Token{RefreshToken: "OLD_REFRESH_TOKEN"})
	id, err := primitive.ObjectIDFromHex("65517e864ac7a12f63fdae33")
	if err != nil {
		t.Fatal(err)
	}

	va1, err := client1.FindVendorAccountById(id)
	if err != nil {
		t.Fatal(err)
	}

	va2, err := client2.FindVendorAccountById(id)
	if err != nil {
		t.Fatal(err)
	}

	tkr1 := client1.NewVendorTokenSource(va1)
	tkr2 := client2.NewVendorTokenSource(va2)

	var tk1 *oauth2.Token
	var tk2 *oauth2.Token

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(tkr oauth2.TokenSource, wg *sync.WaitGroup){
		defer wg.Done()
		var err error
		tk1, err = tkr.Token()
		if err != nil {
			t.Errorf("got err = %v; want none", err)
			return
		}
	}(tkr1, wg)

	go func(tkr oauth2.TokenSource, wg *sync.WaitGroup){
		defer wg.Done()
		var err error
		tk2, err = tkr.Token()
		if err != nil {
			t.Errorf("got err = %v; want none", err)
			return
		}
	}(tkr2,wg)

	wg.Wait()

	if count != 1 {
		t.Fatalf("Count = %d. Should of only hit the endpoint once.", count)
		return
	}

	if tk1.RefreshToken != tk2.RefreshToken && tk1.RefreshToken != "NEW_REFRESH_TOKEN" {
		t.Fatalf("%s != %s != NEW_REFRESH_TOKN", tk1.RefreshToken, tk2.RefreshToken)
		return
	}
}
