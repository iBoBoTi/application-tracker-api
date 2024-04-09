package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config map .env configs to this struct
type Config struct {
	Timezone               string        `mapstructure:"TIMEZONE"`
	DBHost                 string        `mapstructure:"DB_HOST"`
	DBUser                 string        `mapstructure:"DB_USER"`
	DBPassword             string        `mapstructure:"DB_PASSWORD"`
	DBName                 string        `mapstructure:"DB_NAME"`
	DBPort                 string        `mapstructure:"DB_PORT"`
	HttpServerAddress      string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFERESH_TOKEN_DURATION"`
	UNIDOC_LICENSE_API_KEY string        `mapstructure:"UNIDOC_LICENSE_API_KEY"`
	Environment            string        `mapstructure:"ENVIRONMENT"`
}

func Load(p string) (cfg Config, err error) {
	return loader(p, ".env")
}

func LoadWithPath(p string, env string) (cfg Config, err error) {
	return loader(p, env)
}

func loader(p string, env string) (cfg Config, err error) {
	viper.AddConfigPath(p)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
