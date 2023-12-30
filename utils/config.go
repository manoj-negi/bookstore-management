package util

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	PORT        string `mapstructure:"PORT"`
	DB_URI      string `mapstructure:"DB_URI"`
	API_SECRET  string `mapstructure:"API_SECRET"`
	AWS_KEY     string `mapstructure:"AWS_KEY"`
	AWS_SECRET  string `mapstructure:"AWS_SECRET"`
	BUCKET_NAME string `mapstructure:"BUCKET_NAME"`
	BUCKET_URL  string `mapstructure:"BUCKET_URL"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	err = viper.Unmarshal(&config)
	return
}
