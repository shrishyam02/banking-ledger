package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/shrishyam02/banking-ledger/common/logger"

	"github.com/gin-gonic/gin"
)

func basicAuthMiddleware(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		providedUsername, providedPassword, ok := parseBasicAuth(auth)
		if !ok || providedUsername != config.ApiAuth.UserName || providedPassword != config.ApiAuth.Password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		c.Next()
	}
}

func parseBasicAuth(auth string) (string, string, bool) {
	hparts := strings.Split(auth, " ")
	if len(hparts) == 2 {
		decoded, err := base64.StdEncoding.DecodeString(hparts[1])
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error decoding base64 auth")
			return "", "", false
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
	}
	return " ", " ", false
}
