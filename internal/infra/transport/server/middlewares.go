package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"
	"github.com/rs/zerolog"
)

func loggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.Query().Encode()
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		requestSize := c.Request.ContentLength

		midLogger := logger.With().
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ip", ip).
			Str("user_agent", userAgent).
			Str("referer", referer).
			Int64("request_size", requestSize).
			Logger()

		midLogger.Info().Msg("request received")

		c.Next()

		latency := time.Since(now)
		status := c.Writer.Status()
		responseSize := c.Writer.Size()

		midLogger = midLogger.With().
			Int("status", status).
			Dur("latency", latency).
			Int("response_size", responseSize).
			Logger()

		switch {
		case status >= 500:
			midLogger.Error().Msg("request completed")
		case status >= 400:
			midLogger.Warn().Msg("request completed")
		default:
			midLogger.Info().Msg("request completed")
		}
	}
}

// GinToolboxErrorsMiddleware преобразует toolbox-ошибки в HTTP-ответы.
func errorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		lastErr := c.Errors.Last().Err
		if pkgErr, ok := lastErr.(*pkgErr.Error); ok {
			status, payload := pkgErr.ToHTTPResponse()
			c.JSON(status, payload)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    pkgErr.CodeInternal,
			"message": lastErr.Error(),
		})
	}
}

// AuthMiddleware is a middleware for authentication
func authMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			_ = c.Error(ErrAuthError)
			c.Abort()
			return
		}

		if strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")) != token {
			_ = c.Error(ErrAuthError)
			c.Abort()
			return
		}

		c.Next()
	}
}
