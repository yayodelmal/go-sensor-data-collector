package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yayodelmal/go-sensor-data-collector/internal/models"
    "github.com/yayodelmal/go-sensor-data-collector/internal/repository"
)

// SHT31Handler maneja las peticiones HTTP para el sensor SHT31
type SHT31Handler struct {
    repo *repository.SHT31Repository
}

// NewSHT31Handler crea una nueva instancia del handler con el repositorio inyectado
func NewSHT31Handler(repo *repository.SHT31Repository) *SHT31Handler {
    return &SHT31Handler{repo: repo}
}

// SHT31Payload define la estructura JSON que esperamos en POST /sensor/sht31
type SHT31Payload struct {
    Temperature float64 `json:"temperature" binding:"required"`
    Humidity    float64 `json:"humidity" binding:"required"`
}

// RegisterRoutes asocia el handler al router group correspondiente
func (h *SHT31Handler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.POST("/sht31", h.PostReading)
}

// PostReading procesa POST /sensor/sht31
func (h *SHT31Handler) PostReading(c *gin.Context) {
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

    if err := h.repo.Create(c.Request.Context(), &reading); err != nil {
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
