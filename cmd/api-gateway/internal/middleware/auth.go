package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-traffic-analytics/internal/shared/models"
)

func APIKeyAuth(validAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorResponse := models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header required",
			}
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		// Support both "Bearer <key>" and "ApiKey <key>" formats
		var apiKey string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &apiKey)
		if err != nil {
			_, err = fmt.Sscanf(authHeader, "ApiKey %s", &apiKey)
			if err != nil {
				errorResponse := models.ErrorResponse{
					Error:   "Unauthorized",
					Message: "Invalid authorization format",
				}
				c.JSON(http.StatusUnauthorized, errorResponse)
				c.Abort()
				return
			}
		}

		if apiKey != validAPIKey {
			errorResponse := models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid API key",
			}
			c.JSON(http.StatusUnauthorized, errorResponse)
			c.Abort()
			return
		}

		c.Next()
	}
}
