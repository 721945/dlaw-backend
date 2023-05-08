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
	MeiliHost                     string `mapstructure:"MEILI_HOST"`
	MeiliApiKey                   string `mapstructure:"MEILI_MASTER_KEY"`
	GoogleCredPath                string `mapstructure:"GOOGLE_CRED_PATH"`
	SMTPName                      string `mapstructure:"SMTP_NAME"`
	SMTPAddress                   string `mapstructure:"SMTP_ADDRESS"`
	SMTPPassword                  string `mapstructure:"SMTP_PASSWORD"`
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
		MeiliHost:                     viper.GetString("MEILI_HOST"),
		MeiliApiKey:                   viper.GetString("MEILI_MASTER_KEY"),
		GoogleCredPath:                viper.GetString("GOOGLE_CRED_PATH"),
		SMTPName:                      viper.GetString("SMTP_NAME"),
		SMTPAddress:                   viper.GetString("SMTP_ADDRESS"),
		SMTPPassword:                  viper.GetString("SMTP_PASSWORD"),
	}

	//log.Fatal("ERROR 1")

	log.Println("👻 Loaded .env file")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("👹 Can't read .env file")

	}

	log.Println("👻 Unmarshal .env file")

	if err := viper.Unmarshal(&env); err != nil {
		log.Println("👹 Can't loaded: ", err)
	}

	return env
}
