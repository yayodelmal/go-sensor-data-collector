// filepath: internal/db/postgres.go
package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yayodelmal/go-sensor-data-collector/internal/models"
)

var DB *gorm.DB

// Init ahora recibe el DSN
func Init(dsn string) {
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[db] cannot connect: %v", err)
	}
	log.Println("[db] connected")

	//  Auto-migrar los modelos para crear las tablas
	if err := dbConn.AutoMigrate(
		&models.SHT31Reading{},
		&models.DS18B20Reading{},
	); err != nil {
		log.Fatalf("[db] migration failed: %v", err)
	}
	log.Println("[db] migrated tables")

	DB = dbConn
}
