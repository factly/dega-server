package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB - gorm DB
var DB *gorm.DB

// SetupDB is database setuo
func SetupDB() {

	//  TEST ENVIRONMENT SET UP
	IP := os.Getenv("IP_ADDR")
	var err error
	dburl := fmt.Sprintf("postgres://postgres:postgres@" + IP + ":5432/dega?sslmode=disable")
	DB, err = gorm.Open("postgres", dburl)

	// DEPLOYMENT SETUP
	// DSN := os.Getenv("DSN")
	// var err error
	// DB, err = gorm.Open("postgres", DSN)

	if err != nil {
		log.Fatal(err)
	}

	// Query log
	DB.LogMode(true)

	fmt.Println("connected to database ...")
}
