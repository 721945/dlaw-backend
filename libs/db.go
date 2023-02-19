package libs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	DB *gorm.DB
}

func NewDatabase(env Env, logger *Logger) Database {
	var url string

	url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=verify-full TimeZone=Asia/Bangkok", env.DBHost, env.DBUser, env.DBPassword, env.DBName)
	//url = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok", env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		logger.Fatal("👹 Can't connect to database: ", err)
	}

	return Database{DB: db}
}

//func (db *Database) Migrates(x uint) {
//	db.DB.AutoMigrate(&models.User{})
//}
