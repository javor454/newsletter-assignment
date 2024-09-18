package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/logger"
)

func LoggingMiddleware(lg logger.Logger, blacklist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, blacklistedPath := range blacklist {
			if strings.Contains(c.Request.URL.Path, blacklistedPath) {
				c.Next()
				return
			}
		}

		start := time.Now()

		var requestBody []byte
		var err error
		if c.Request.Body != nil {
			requestBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				panic("failed to read request body: " + err.Error())
			}
		}
		// Restore the request body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Create a response writer that captures the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		duration := time.Since(start)
		meta := map[string]interface{}{
			"duration": duration.String(),
		}
		if string(requestBody) != "" {
			meta["request_body"] = string(requestBody)
		}
		if blw.body.String() != "" {
			meta["response_body"] = blw.body.String()
		}

		lg.WithFields(meta).Infof("Request / Response: %s %s (%d)", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
