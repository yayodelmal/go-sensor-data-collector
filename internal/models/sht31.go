package models

import "time"

// SHT31Reading representa un registro de temperatura + humedad
// en la tabla `sht31_readings`. El campo TS tomará DEFAULT now() en la DB.
type SHT31Reading struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Temperature float64   `gorm:"not null" json:"temperature"`
    Humidity    float64   `gorm:"not null" json:"humidity"`
    TS          time.Time `gorm:"not null;default:now()" json:"datetime"`
}

// TableName fuerza a GORM a usar la tabla “sht31_readings”
func (SHT31Reading) TableName() string {
    return "sht31_readings"
}
