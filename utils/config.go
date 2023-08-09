package utils

import "github.com/spf13/viper"

type Config struct {
	DbDriver   string `mapstructure:"DB_DRIVER"`
	DbSource   string `mapstructure:"DB_SOURCE"`
	ServerAddr string `mapstructure:"SERVER_ADDR"`
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
