package libs

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Environment string `mapstructure:"ENVIRONMENT"`
	//DSN        string `mapstructure:"DSN"`
}

func NewEnv() Env {
	env := Env{}
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("👹 Can't read .env file")
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal("👹 Can't loaded: ", err)
	}

	return env
}
