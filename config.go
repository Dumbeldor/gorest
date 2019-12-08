package gorest

import (
	"fmt"

	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() {
	setDefault()

	var env = ""
	getopt.FlagLong(&env, "env", 'E', "Set application environment")
	getopt.Parse()

	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		log.Errorln("Unable to read config file", err)
		return
	}
	log.Infoln(fmt.Sprintf("Config '%s' loaded", env))
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func setDefault() {
	viper.SetDefault("port", "8080")
}
