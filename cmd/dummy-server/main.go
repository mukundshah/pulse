package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	r := gin.Default()

	r.Any("/*path", func(c *gin.Context) {
		simulateRandomBehavior(c)
	})

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// simulateRandomBehavior simulates various random HTTP behaviors for testing
func simulateRandomBehavior(c *gin.Context) {
	// 1% chance of timeout (hang indefinitely)
	if rand.Float32() < 0.01 {
		select {} // Block forever
	}

	// Random delay (0-5 seconds)
	delay := rand.Float64() * 5
	if delay > 0 {
		time.Sleep(time.Duration(delay * float64(time.Second)))
	}

	// Random behavior selection
	behavior := rand.Intn(100)

	switch {
	// 60% chance of success (200, 201)
	case behavior < 60:
		statusCode := []int{http.StatusOK, http.StatusCreated}[rand.Intn(2)]
		respondWithStatus(c, statusCode)

	// 5% chance of client errors (400, 401, 403, 404)
	case behavior < 65:
		statusCode := []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}[rand.Intn(4)]
		respondWithStatus(c, statusCode)

	// 20% chance of server errors (500, 502, 503, 504)
	case behavior < 85:
		statusCode := []int{
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		}[rand.Intn(4)]
		respondWithStatus(c, statusCode)

	// 5% chance of rate limiting (429)
	case behavior < 90:
		c.Header("Retry-After", fmt.Sprintf("%d", rand.Intn(60)+1))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       "Rate limit exceeded",
			"retry_after": rand.Intn(60) + 1,
		})

	// 3% chance of redirect
	case behavior < 93:
		redirectCodes := []int{
			http.StatusMovedPermanently,
			http.StatusFound,
			http.StatusSeeOther,
			http.StatusTemporaryRedirect,
		}
		code := redirectCodes[rand.Intn(len(redirectCodes))]
		c.Redirect(code, "/")

	// 2% chance of connection reset
	case behavior < 95:
		if conn, ok := c.Writer.(http.Hijacker); ok {
			if raw, _, err := conn.Hijack(); err == nil {
				raw.Close()
				return
			}
		}
		c.String(http.StatusOK, "Connection reset")

	// 2% chance of malformed response
	case behavior < 97:
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.WriteString("This is not valid JSON{")

	// 3% chance of chunked response
	default:
		c.Header("Transfer-Encoding", "chunked")
		c.Writer.WriteHeader(http.StatusOK)
		chunks := rand.Intn(5) + 1
		for i := 0; i < chunks; i++ {
			c.Writer.WriteString(fmt.Sprintf("Chunk %d\n", i+1))
			c.Writer.Flush()
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
	}
}

// respondWithStatus responds with a random content type based on status code
func respondWithStatus(c *gin.Context, statusCode int) {
	contentType := rand.Intn(3)

	switch contentType {
	case 0: // JSON
		c.JSON(statusCode, gin.H{
			"status":    statusCode,
			"message":   http.StatusText(statusCode),
			"timestamp": time.Now().Unix(),
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
		})
	case 1: // Plain text
		c.String(statusCode, "Status: %d\nMessage: %s\nPath: %s\nMethod: %s",
			statusCode, http.StatusText(statusCode), c.Request.URL.Path, c.Request.Method)
	case 2: // HTML
		html := fmt.Sprintf(
			"<html><body><h1>Status %d</h1><p>%s</p><p>Path: %s</p><p>Method: %s</p></body></html>",
			statusCode, http.StatusText(statusCode), c.Request.URL.Path, c.Request.Method)
		c.Data(statusCode, "text/html", []byte(html))
	}
}
