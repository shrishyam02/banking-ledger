package server

import (
	"common/logger"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func basicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		providedUsername, providedPassword, ok := parseBasicAuth(auth)

		username := os.Getenv("API_AUTH_USERNAME")
		password := os.Getenv("API_AUTH_PASSWORD")
		if !ok || providedUsername != username || providedPassword != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		c.Next()
	}
}

func parseBasicAuth(auth string) (string, string, bool) {
	if len(auth) > 6 && strings.ToUpper(auth[:6]) == "BASIC" {
		b64 := auth[6:]
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error decoding base64 auth")
			return "", "", false
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
	}
	return "", "", false
}
