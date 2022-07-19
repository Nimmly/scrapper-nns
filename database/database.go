package database

import (
	"fmt"
	"log"

	models "github.com/Nimmly/scrapper-nns/models"
	util "github.com/Nimmly/scrapper-nns/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Db *gorm.DB
}

var Database DB

func ConnectDb() {
	config, err := util.ReadPostgresConfig("config_files/postgres_db.yaml")
	if err != nil {
		log.Println("Failed to read YAML file!")
	}
	//============================================================================
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(postgres.Open(dbinfo))
	if err != nil {
		log.Fatal("Postgres instance terminated... Check your connection!")
	}
	log.Println("Connected to the database")
	db.Logger = logger.Default.LogMode(logger.Info)
	//============================================================================
	db.AutoMigrate(&models.Apartment{}, &models.City{}, &models.Image{})
	Database = DB{Db: db}
}
