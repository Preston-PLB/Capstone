package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Mongo     *MongoConfig `mapstructure:"mongo"`
	JwtSecret string       `mapstructure:"jwt_secret"`
	Env       string       `mapstructure:"env"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	EntDb  string `mapstructure:"ent_db"`
	EntCol string `mapstructure:"ent_col"`
}

var cfg *config

func Init() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
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
