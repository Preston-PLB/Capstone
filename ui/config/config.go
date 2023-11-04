package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Mongo     *MongoConfig             `mapstructure:"mongo"`
	Vendors   map[string]*VendorConfig `mapstructure:"vendors"`
	JwtSecret string                   `mapstructure:"jwt_secret"`
	Env       string                   `mapstructure:"env"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	EntDb  string `mapstructure:"ent_db"`
	EntCol string `mapstructure:"ent_col"`
}

type VendorConfig struct {
	ClientId      string   `mapstructure:"client_id"`
	ClientSecret  string   `mapstructure:"client_secret"`
	Scopes        []string `mapstructure:"scopes"`
	AuthUri       string   `mapstructure:"auth_uri"`
	TokenUri      string   `mapstructure:"token_uri"`
	RefreshEncode string   `mapstructure:"refresh_encode"`
	scope         string
}

func (pco *VendorConfig) Scope() string {
	if pco.scope == "" {
		pco.scope = strings.Join(pco.Scopes, " ")
	}
	return pco.scope
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

	fmt.Printf("%v\n", cfg)
	for key, value := range cfg.Vendors {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func Config() *config {
	return cfg
}
