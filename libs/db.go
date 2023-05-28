package libs

import (
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// Config : Unused here
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type Database struct {
	DB    *gorm.DB
	Meili *meilisearch.Client
}

func NewDatabase(env Env, myLogger *Logger) Database {
	var url string

	url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=verify-full TimeZone=Asia/Bangkok", env.DBHost, env.DBUser, env.DBPassword, env.DBName)
	//url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		myLogger.Fatal("👹 Can't connect to database: ", err)
	}

	myLogger.Info("👻 Connected to database - 1")

	// Uncomment if want to Migrate the schema
	//err = db.AutoMigrate(
	//	&models.Case{},
	//	&models.Action{},
	//	&models.Permission{},
	//	&models.ActionLog{},
	//	&models.Appointment{},
	//	&models.Email{},
	//	&models.CasePermission{},
	//	&models.CasePermissionLog{},
	//	&models.CaseUsedLog{},
	//	&models.FileType{},
	//	&models.File{},
	//	&models.Folder{},
	//	&models.Tag{},
	//	&models.User{},
	//	&models.FileViewLog{},
	//)
	//&models.FileVersion{},
	//	&models.FileUrl{},

	if err != nil {
		myLogger.Fatal("👹 Can't migrate database: ", err)
	}

	meili := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   env.MeiliHost,
		APIKey: env.MeiliApiKey,
	})

	return Database{DB: db, Meili: meili}
}

//func (db *Database) Migrates(x uint) {
//	db.DB.AutoMigrate(&models.User{})
//}
