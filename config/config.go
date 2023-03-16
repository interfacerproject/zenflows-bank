package config

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	ZenflowsUrl  string `mapstructure:"ZENFLOWS_URL"`
	ZenflowsSk   string `mapstructure:"ZENFLOWS_SK"`
	ZenflowsUser string `mapstructure:"ZENFLOWS_USER"`
	TTHost       string `mapstructure:"TT_HOST"`
	TTUser       string `mapstructure:"TT_USER"`
	TTPass       string `mapstructure:"TT_PASS"`
	Fabcoin      string `mapstructure:"FABCOIN"`
	EthereumUrl      string `mapstructure:"ETHEREUM_URL"`
	EthereumSk      string `mapstructure:"ETHEREUM_SK"`
}

var Config *EnvConfig

func Init() {
	var err error
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err.Error())
	}
}
