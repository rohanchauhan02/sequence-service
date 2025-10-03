package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	ImmutableConfig interface {
		GetPort() string
	}

	config struct {
		Port string `mapstructure:"PORT"`
	}
)

var (
	once sync.Once
	conf *config
)

func NewImmutableConfig() ImmutableConfig {
	once.Do(func() {
		v := viper.New()
		appEnv, exists := os.LookupEnv("APP_ENV")
		configName := "app.config.local"
		if exists {
			switch appEnv {
			case "development":
				configName = "app.config.dev"
			case "production":
				configName = "app.config.prod"
			}
		}

		slog.Debug("Config loaded", slog.String("ConfigName", configName), slog.String("Level", "warn"))

		v.SetConfigName("configs/" + configName)
		v.AddConfigPath(".")

		v.SetEnvPrefix("GO_SEQUENCE")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			panic(err.Error())
		}

		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		err := v.Unmarshal(&conf)
		if err != nil {
			panic(err.Error())
		}
	})
	return conf
}

func (c *config) GetPort() string {
	return c.Port
}
