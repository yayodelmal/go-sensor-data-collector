package models

import "time"

// DS18B20Reading representa un registro de temperatura
// en la tabla `ds18b20_readings`. El campo TS tomará DEFAULT now().
type DS18B20Reading struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Temperature float64   `gorm:"not null" json:"temperature"`
    TS          time.Time `gorm:"not null;default:now()" json:"datetime"`
}

// TableName fuerza a GORM a usar la tabla “ds18b20_readings”
func (DS18B20Reading) TableName() string {
    return "ds18b20_readings"
}
