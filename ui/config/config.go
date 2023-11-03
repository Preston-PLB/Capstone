package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Mongo         *MongoConfig   `mapstructure:"mongo"`
	YoutubeConfig *YoutubeConfig `mapstructure:"youtube"`
	PcoConfig     *PcoConfig     `mapstructure:"pco"`
	JwtSecret     string         `mapstructure:"jwt_secret"`
	Env           string         `mapstructure:"env"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	EntDb  string `mapstructure:"ent_db"`
	EntCol string `mapstructure:"ent_col"`
}

type YoutubeConfig struct {
	ClientId     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	Scopes       []string `mapstructure:"scopes"`
	AuthUri      string   `mapstructure:"auth_uri"`
	TokenUri     string   `mapstructure:"token_uri"`
	scope        string
}

func (yt *YoutubeConfig) Scope() string {
	if yt.scope == "" {
		for i, str := range yt.Scopes {
			yt.Scopes[i] = fmt.Sprintf("https://www.googleapis.com%s", str)
		}
		yt.scope = strings.Join(yt.Scopes, " ")
	}
	return yt.scope
}

type PcoConfig struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	AuthUri      string `mapstructure:"auth_uri"`
	TokenUri     string `mapstructure:"token_uri"`
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
