// filepath: internal/repository/sensor_repository.go
package repository

import "context"

// SensorRepository define las operaciones comunes para repositorios de sensores.
// Usa generics (T) para permitir diferentes tipos de lecturas (SHT31Reading, DS18B20Reading, etc.)
type SensorRepository[T any] interface {
	// Create guarda una nueva lectura en la base de datos
	Create(ctx context.Context, reading *T) error
}
