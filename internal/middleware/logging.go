package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()

        // Process request
        c.Next()

        endTime := time.Now()
        latency := endTime.Sub(startTime)

        // Log request details
        log.WithFields(log.Fields{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "latency_ms": latency.Milliseconds(),
            "client_ip":  c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        }).Info("HTTP request processed")
    }
}