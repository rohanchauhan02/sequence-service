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
		GetDBConf() DB
		GetKafkaConf() Kafka
	}

	config struct {
		Port  string `mapstructure:"PORT"`
		DB    DB     `mapstructure:"DB"`
		Kafka Kafka  `mapstructure:"KAFKA"`
	}
	DB struct {
		Host             string `mapstructure:"HOST"`
		Port             int    `mapstructure:"PORT"`
		Name             string `mapstructure:"NAME"`
		User             string `mapstructure:"USER"`
		Password         string `mapstructure:"PASSWORD"`
		MaxIdleConns     int    `mapstructure:"MAX_IDLE_CONNS"`
		MaxOpenConns     int    `mapstructure:"MAX_OPEN_CONNS"`
		MaxLifetimeConns int    `mapstructure:"MAX_LIFETIME_CONNS"`
		SSLMode          string `mapstructure:"SSL_MODE"`
	}
	Kafka struct {
		Broker string `mapstructure:"BROKERS"`
		Topics Topic    `mapstructure:"TOPICS"`
	}

	Topic struct {
		EmailJobs      string `mapstructure:"EMAIL_JOBS"`
		FollowupEvents string `mapstructure:"FOLLOWUP_EVENTS"`
		EmailRetries   string `mapstructure:"EMAIL_RETRIES"`
		EmailEvents    string `mapstructure:"EMAIL_EVENTS"`
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
			default:
				configName = "app.config.dev"
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

func (im *config) GetDBConf() DB {
	return im.DB
}

func (im *config) GetKafkaConf() Kafka {
	return im.Kafka
}
