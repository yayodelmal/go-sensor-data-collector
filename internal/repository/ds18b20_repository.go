// filepath: internal/repository/ds18b20_repository.go
package repository

import (
	"context"

	"github.com/yayodelmal/go-sensor-data-collector/internal/models"
	"gorm.io/gorm"
)

// DS18B20Repository implementa SensorRepository para lecturas del sensor DS18B20
type DS18B20Repository struct {
	db *gorm.DB
}

// NewDS18B20Repository crea una nueva instancia del repositorio DS18B20
func NewDS18B20Repository(db *gorm.DB) *DS18B20Repository {
	return &DS18B20Repository{db: db}
}

// Create guarda una nueva lectura DS18B20 en la base de datos
func (r *DS18B20Repository) Create(ctx context.Context, reading *models.DS18B20Reading) error {
	return r.db.WithContext(ctx).Create(reading).Error
}
