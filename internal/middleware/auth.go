package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el header "Authorization: Bearer <token>"
// contra el valor de la variable de entorno API_TOKEN. Si no coincide,
// responde 401 Unauthorized y aborta la petición.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        expectedToken := os.Getenv("API_TOKEN")
        if expectedToken == "" {
            // Si no se definió API_TOKEN, no aplicamos autenticación.
            c.Next()
            return
        }

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
            return
        }

        receivedToken := strings.TrimPrefix(authHeader, "Bearer ")
        if receivedToken != expectedToken {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // Token válido, continuar con el handler
        c.Next()
    }
}
