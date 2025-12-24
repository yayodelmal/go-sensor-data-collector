// filepath: internal/db/postgres.go
package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yayodelmal/go-sensor-data-collector/internal/models"
)

var DB *gorm.DB

// Init ahora recibe el DSN
func Init(dsn string) {
	// Configuración de GORM con PrepareStmt para mejor performance
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // Usa prepared statements (mejor performance)
	})
	if err != nil {
		log.Fatalf("[db] cannot connect: %v", err)
	}
	log.Println("[db] connected")

	// Configurar connection pool
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatalf("[db] failed to get sql.DB: %v", err)
	}

	// SetMaxIdleConns establece el número máximo de conexiones idle
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns establece el número máximo de conexiones abiertas
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime establece el tiempo máximo que una conexión puede ser reutilizada
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime establece el tiempo máximo que una conexión puede estar idle
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Println("[db] connection pool configured")

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
