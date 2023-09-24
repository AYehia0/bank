package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver                   string        `mapstructure:"DB_DRIVER"`
	DbSource                   string        `mapstructure:"DB_SOURCE"`
	ServerAddr                 string        `mapstructure:"SERVER_ADDR"`
	TokenKey                   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenExpireDuration        time.Duration `mapstructure:"TOKEN_EXPIRE_TIME"`
	TokenRefreshExpireDuration time.Duration `mapstructure:"TOKEN_REFRESH_EXPIRE_TIME"`
}

func ConfigStore(configPath, configName, configType string) (config Config, err error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// also read from env vars
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}
