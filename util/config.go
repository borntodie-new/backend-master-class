package util

import (
	"github.com/spf13/viper"
	"time"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME_QQ"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS_QQ"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD_QQ"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	// set config file path
	viper.AddConfigPath(path)
	// set config file name
	viper.SetConfigName("app")
	// set config file type
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
