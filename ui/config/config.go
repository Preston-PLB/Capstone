package config

import (
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type config struct {
	Mongo       *MongoConfig             `mapstructure:"mongo"`
	Vendors     map[string]*VendorConfig `mapstructure:"vendors"`
	JwtSecret   string                   `mapstructure:"jwt_secret"`
	Env         string                   `mapstructure:"env"`
	AppSettings *AppSettings             `mapstructure:"app_settings"`
}

type AppSettings struct {
	WebhookServiceUrl  string `mapstructure:"webhook_service_url"`
	FrontendServiceUrl string `mapstructure:"frontend_service_url"`
}

type MongoConfig struct {
	Uri     string `mapstructure:"uri"`
	EntDb   string `mapstructure:"ent_db"`
	EntCol  string `mapstructure:"ent_col"`
	LockDb  string `mapstructure:"lock_db"`
	LockCol string `mapstructure:"lock_col"`
}

type VendorConfig struct {
	ClientId      string   `mapstructure:"client_id"`
	ClientSecret  string   `mapstructure:"client_secret"`
	Scopes        []string `mapstructure:"scopes"`
	AuthUri       string   `mapstructure:"auth_uri"`
	TokenUri      string   `mapstructure:"token_uri"`
	RefreshEncode string   `mapstructure:"refresh_encode"`
	WebhookSecret string   `mapstructure:"webhook_secret"`
	scope         string
}

func (vendor *VendorConfig) Scope() string {
	if vendor.scope == "" {
		vendor.scope = strings.Join(vendor.Scopes, " ")
	}
	return vendor.scope
}

func (vendor *VendorConfig) OauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     vendor.ClientId,
		ClientSecret: vendor.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   vendor.AuthUri,
			TokenURL:  vendor.TokenUri,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "",
		Scopes:      vendor.Scopes,
	}
}

var cfg *config

func Init() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/capstone") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg = &config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
}

func Config() *config {
	return cfg
}
