package libs

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	ServerPort                    string `mapstructure:"SERVER_PORT"`
	LogOutput                     string `mapstructure:"LOG_OUTPUT"`
	DBHost                        string `mapstructure:"DB_HOST"`
	DBPort                        string `mapstructure:"DB_PORT"`
	DBUser                        string `mapstructure:"DB_USER"`
	DBPassword                    string `mapstructure:"DB_PASSWORD"`
	DBName                        string `mapstructure:"DB_NAME"`
	JWTSecret                     string `mapstructure:"JWT_SECRET"`
	Environment                   string `mapstructure:"ENVIRONMENT"`
	Bucket                        string `mapstructure:"BUCKET"`
	GoogleCloudStoragePrivateKey  string `mapstructure:"GOOGLE_CLOUD_STORAGE_PRIVATE_KEY"`
	GoogleCloudStorageClientEmail string `mapstructure:"GOOGLE_CLOUD_STORAGE_CLIENT_EMAIL"`
}

func NewEnv() Env {
	viper.AutomaticEnv()

	env := Env{
		ServerPort:                    viper.GetString("SERVER_PORT"),
		LogOutput:                     viper.GetString("LOG_OUTPUT"),
		DBHost:                        viper.GetString("DB_HOST"),
		DBPort:                        viper.GetString("DB_PORT"),
		DBUser:                        viper.GetString("DB_USER"),
		DBPassword:                    viper.GetString("DB_PASSWORD"),
		DBName:                        viper.GetString("DB_NAME"),
		JWTSecret:                     viper.GetString("JWT_SECRET"),
		Environment:                   viper.GetString("ENVIRONMENT"),
		Bucket:                        viper.GetString("BUCKET"),
		GoogleCloudStoragePrivateKey:  viper.GetString("GOOGLE_CLOUD_STORAGE_PRIVATE_KEY"),
		GoogleCloudStorageClientEmail: viper.GetString("GOOGLE_CLOUD_STORAGE_CLIENT_EMAIL"),
	}

	//log.Fatal("ERROR 1")

	viper.SetConfigType("env")

	viper.SetConfigFile(".env")

	log.Println(viper.GetString("DB_HOST"))

	log.Println("👻 Loaded .env file")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("👹 Can't read .env file")

	}

	log.Println("👻 Unmarshal .env file")

	if err := viper.Unmarshal(&env); err != nil {
		log.Println("👹 Can't loaded: ", err)
	}

	log.Println(env)

	return env
}
