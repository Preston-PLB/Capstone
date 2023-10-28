package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Mongo     *MongoConfig `json:"mogno"`
	JwtSecret string       `json:"jwt_secret"`
	Env       string       `json:"env"`
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	EntDb    string `json:"entity_db"`
	EntCol   string `json:"entity_col"`
}

var cfg *config

func init() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/capstone/") // path to look for the config file in

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
