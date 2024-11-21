package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}

	// Enable uuid-ossp extension
	if err := DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		fmt.Println("Error enabling uuid-ossp extension")
		panic(err)
	}

	err = DBConn.AutoMigrate(&User{}, &SearchSetting{}, &CrawledUrl{}, &SearchIndex{})
	if err != nil {
		fmt.Println("Error migrating database")
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConn
}
