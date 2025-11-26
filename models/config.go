package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type EnvConfig struct {
	ServerPort                 string `mapstructure:"SERVER_PORT"`
	ServerAddr                 string `mapstructure:"SERVER_ADDR"`
	MySQLHost                  string `mapstructure:"MYSQL_HOST"`
	MySQLPort                  string `mapstructure:"MYSQL_PORT"`
	MySQLUser                  string `mapstructure:"MYSQL_USER"`
	MySQLPassword              string `mapstructure:"MYSQL_PASSWORD"`
	MySQLDatabase              string `mapstructure:"MYSQL_DATABASE"`
	MySQLCharset               string `mapstructure:"MYSQL_CHARSET"`
	UseRedis                   bool   `mapstructure:"USE_REDIS"`
	RedisDefaultAddr           string `mapstructure:"REDIS_DEFAULT_ADDR"`
	JWTSecretKey               string `mapstructure:"JWT_SECRET"`
	JWTAccessExpirationMinutes int    `mapstructure:"JWT_ACCESS_EXPIRATION_MINUTES"`
	JWTRefreshExpirationDays   int    `mapstructure:"JWT_REFRESH_EXPIRATION_DAYS"`
	Mode                       string `mapstructure:"MODE"` // Added closing quotation mark
}

func (config *EnvConfig) Validate() error {
	return validation.ValidateStruct(config,
		validation.Field(&config.ServerPort, is.Port),
		validation.Field(&config.ServerAddr, validation.Required),

		validation.Field(&config.MySQLHost, validation.Required),
		validation.Field(&config.MySQLPort, validation.Required),
		validation.Field(&config.MySQLUser, validation.Required),
		validation.Field(&config.MySQLPassword, validation.Required),
		validation.Field(&config.MySQLDatabase, validation.Required),
		validation.Field(&config.MySQLCharset, validation.In("utf8", "utf8mb4")),

		validation.Field(&config.UseRedis, validation.In(true, false)),
		validation.Field(&config.RedisDefaultAddr),

		validation.Field(&config.JWTSecretKey, validation.Required),
		validation.Field(&config.JWTAccessExpirationMinutes, validation.Required),
		validation.Field(&config.JWTRefreshExpirationDays, validation.Required),

		validation.Field(&config.Mode, validation.In("debug", "release")),
	)
}
