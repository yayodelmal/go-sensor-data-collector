// filepath: internal/repository/sht31_repository.go
package repository

import (
	"context"

	"github.com/yayodelmal/go-sensor-data-collector/internal/models"
	"gorm.io/gorm"
)

// SHT31Repository implementa SensorRepository para lecturas del sensor SHT31
type SHT31Repository struct {
	db *gorm.DB
}

// NewSHT31Repository crea una nueva instancia del repositorio SHT31
func NewSHT31Repository(db *gorm.DB) *SHT31Repository {
	return &SHT31Repository{db: db}
}

// Create guarda una nueva lectura SHT31 en la base de datos
func (r *SHT31Repository) Create(ctx context.Context, reading *models.SHT31Reading) error {
	return r.db.WithContext(ctx).Create(reading).Error
}
