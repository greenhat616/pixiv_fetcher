package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

// initConfig 用于初始化配置
func initConfig() {
	setConfigDefaults()
	// Parse env config
	viper.SetEnvPrefix("pixiv_fetcher") // like: PIXIV_FETCHER_PORT=8000
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default viper information
	viper.SetConfigName("config")
	viper.SetConfigType("toml") // Toml is the best!

	// Parse path etc > home > localPath
	viper.AddConfigPath("/etc/.pixiv_fetcher")
	viper.AddConfigPath("$HOME/.pixiv_fetcher")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("[config] Fatal error while reading config file: %s \n", err)
	}
}

// setConfigDefaults 用于设置初始值
func setConfigDefaults() {
	viper.SetDefault("server.port", 8000)
}
