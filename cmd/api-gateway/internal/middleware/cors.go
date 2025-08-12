package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS es un middleware que maneja los headers de CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el origen de la solicitud
		origin := c.Request.Header.Get("Origin")

		// Lista de orígenes permitidos (en producción, haz esto configurable)
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
			// Agrega tus dominios permitidos aquí
		}

		// Verificar si el origen está permitido
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed || allowed == "*" {
				isAllowed = true
				break
			}
		}

		if isAllowed || len(allowedOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Headers permitidos
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-API-Key",
		}, ", "))

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Manejar preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
