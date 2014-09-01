package remailer

import (
	"log"
	"os"
	"time"
	"github.com/gin-gonic/gin"
)

var (
	green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

func Logger() gin.HandlerFunc {
	stdlogger := log.New(os.Stdout, "", 0)
	//errlogger := log.New(os.Stderr, "", 0)

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		c.Set("email", "")
		c.Set("body", "")

		// Process request
		c.Next()

		// save the IP of the requester
		requester := c.Request.Header.Get("X-Real-IP")
		// if the requester-header is empty, check the forwarded-header
		if len(requester) == 0 {
			requester = c.Request.Header.Get("X-Forwarded-For")
		}
		// if the requester is still empty, use the hard-coded address from the socket
		if len(requester) == 0 {
			requester = c.Request.RemoteAddr
		}
		email := c.MustGet("email").(string)
		body := c.MustGet("body").(string)

		var color string
		code := c.Writer.Status()
		switch {
		case code >= 200 && code <= 299:
			color = green
		case code >= 300 && code <= 399:
			color = white
		case code >= 400 && code <= 499:
			color = red
		default:
			color = red
		}
		end := time.Now()
		latency := end.Sub(start)
		stdlogger.Printf("[GIN] %v |%s %3d %s| %12v | %s %4s %s\n%s | %v | %v",
			end.Format("2006/01/02 - 15:04:05"),
			color, code, reset,
			latency,
			requester,
			c.Request.Method, c.Request.URL.Path,
			c.Errors.String(), email, body)
	}
}