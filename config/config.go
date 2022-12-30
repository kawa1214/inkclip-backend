package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config contains all the configuration for the application
// The values read by viper from the config file are mapped to this struct
type Config struct {
	Env                  string        `mapstructure:"ENV"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecretKey       string        `mapstructure:"TOKEN_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MailHostname         string        `mapstructure:"MAIL_HOSTNAME"`
	MailPort             int           `mapstructure:"MAIL_PORT"`
	MailUsername         string        `mapstructure:"MAIL_USERNAME"`
	MailPassword         string        `mapstructure:"MAIL_PASSWORD"`
	FrontURL             string        `mapstructure:"FRONT_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("ENV", "")
	viper.SetDefault("DB_DRIVER", "")
	viper.SetDefault("DB_SOURCE", "")
	viper.SetDefault("MIGRATION_URL", "")
	viper.SetDefault("SERVER_ADDRESS", "")
	viper.SetDefault("TOKEN_SECRET_KEY", "")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "")
	viper.SetDefault("MAIL_HOSTNAME", "")
	viper.SetDefault("MAIL_PORT", "")
	viper.SetDefault("MAIL_USERNAME", "")
	viper.SetDefault("MAIL_PASSWORD", "")
	viper.SetDefault("FRONT_URL", "")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read config: ", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("cannot Unmarshal config: ", err)
		return
	}

	return
}
