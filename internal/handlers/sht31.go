package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yayodelmal/go-sensor-data-collector/internal/db"
    "github.com/yayodelmal/go-sensor-data-collector/internal/models"
)

// SHT31Payload define la estructura JSON que esperamos en POST /sensor/sht31
type SHT31Payload struct {
    Temperature float64 `json:"temperature" binding:"required"`
    Humidity    float64 `json:"humidity" binding:"required"`
}

// RegisterSHT31Routes asocia el handler al router group correspondiente.
func RegisterSHT31Routes(rg *gin.RouterGroup) {
    rg.POST("/sht31", postSHT31)
}

// postSHT31 procesa POST /sensor/sht31
func postSHT31(c *gin.Context) {
    var payload SHT31Payload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":  "invalid JSON payload",
            "detail": err.Error(),
        })
        return
    }

    reading := models.SHT31Reading{
        Temperature: payload.Temperature,
        Humidity:    payload.Humidity,
        // TS se omite: la DB usará DEFAULT now()
    }

    if err := db.DB.Create(&reading).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "could not save SHT31 reading",
            "detail": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "SHT31 reading saved",
        "record":  reading,
    })
}
