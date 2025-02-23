package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"common/logger"

	"github.com/gin-gonic/gin"
)

// Config holds the server configuration.
type Config struct {
	Port        string
	ServiceName string
}

//HandlerRegistrationFunc ...
type HandlerRegistrationFunc func(router *gin.RouterGroup)

//RunServer ...
func RunServer(ctx context.Context, config Config, registerHandlers HandlerRegistrationFunc) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(RequestLogger())

	apiGroup := router.Group("/api/v1")
	apiGroup.Use(basicAuthMiddleware())
	registerHandlers(apiGroup)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msgf("Failed to start %s server", config.ServiceName)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info().Msgf("Shutting down %s server...", config.ServiceName)

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msgf("%s server shutdown failed", config.ServiceName)
	}

	logger.Log.Info().Msgf("%s server gracefully stopped.", config.ServiceName)
}

// RequestLogger - logging middleware.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Log.Info().
			Int("status", statusCode).
			Str("method", method).
			Str("path", path).
			Dur("latency", latency).
			Msg("Request processed")
	}
}
