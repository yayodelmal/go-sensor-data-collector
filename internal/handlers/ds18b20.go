package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yayodelmal/go-sensor-data-collector/internal/models"
    "github.com/yayodelmal/go-sensor-data-collector/internal/repository"
)

// DS18B20Handler maneja las peticiones HTTP para el sensor DS18B20
type DS18B20Handler struct {
    repo *repository.DS18B20Repository
}

// NewDS18B20Handler crea una nueva instancia del handler con el repositorio inyectado
func NewDS18B20Handler(repo *repository.DS18B20Repository) *DS18B20Handler {
    return &DS18B20Handler{repo: repo}
}

// DS18B20Payload define la estructura JSON que esperamos en POST /sensor/ds18b20
type DS18B20Payload struct {
    Temperature float64 `json:"temperature" binding:"required"`
}

// RegisterRoutes asocia el handler al router group correspondiente
func (h *DS18B20Handler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.POST("/ds18b20", h.PostReading)
}

// PostReading procesa POST /sensor/ds18b20
func (h *DS18B20Handler) PostReading(c *gin.Context) {
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

    if err := h.repo.Create(c.Request.Context(), &reading); err != nil {
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
