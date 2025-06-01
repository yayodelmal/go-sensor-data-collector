package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yayodelmal/go-sensor-data-collector/internal/db"
    "github.com/yayodelmal/go-sensor-data-collector/internal/models"
)

// DS18B20Payload define la estructura JSON que esperamos en POST /sensor/ds18b20
type DS18B20Payload struct {
    Temperature float64 `json:"temperature" binding:"required"`
}

// RegisterDS18B20Routes asocia el handler al router group correspondiente.
func RegisterDS18B20Routes(rg *gin.RouterGroup) {
    rg.POST("/ds18b20", postDS18B20)
}

// postDS18B20 procesa POST /sensor/ds18b20
func postDS18B20(c *gin.Context) {
    var payload DS18B20Payload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":  "invalid JSON payload",
            "detail": err.Error(),
        })
        return
    }

    reading := models.DS18B20Reading{
        Temperature: payload.Temperature,
        // TS se omite: la DB usará DEFAULT now()
    }

    if err := db.DB.Create(&reading).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "could not save DS18B20 reading",
            "detail": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "DS18B20 reading saved",
        "record":  reading,
    })
}
