package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort            string
	Postgres            PostgresConfig
	Smtp                Smtp
	RedisAddr           string
	AuthSecretKey       string
	AuthHeaderKey       string
	AuthPayloadKey      string
	AccessTokenDuration time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Smtp struct {
	Sender   string
	Password string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env") // load .env file if it exists

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		HttpPort: conf.GetString("HTTP_PORT"),
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
		Smtp: Smtp{
			Sender:   conf.GetString("SMTP_SENDER"),
			Password: conf.GetString("SMTP_PASSWORD"),
		},
		RedisAddr:      conf.GetString("REDIS_ADDR"),
		AuthSecretKey:  conf.GetString("AUTH_SECRET_KEY"),
		AuthHeaderKey:  conf.GetString("AUTHORIZATION_HEADER_KEY"),
		AuthPayloadKey: conf.GetString("AUTHORIZATION_PAYLOAD_KEY"),
		AccessTokenDuration: conf.GetDuration("ACCESS_TOKEN_DURATION"),
	}

	return cfg
}
