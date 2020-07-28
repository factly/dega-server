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
	TESTDSN := os.Getenv("TESTDSN")
	var err error
	DB, err = gorm.Open("postgres", TESTDSN)

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
