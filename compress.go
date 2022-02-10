package compress

import (
	"compress/flate"
	"compress/gzip"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncodingsStr := c.GetHeader("accept-encoding")
		if len(acceptEncodingsStr) == 0 {
			return
		}
		acceptEncodings := strings.Split(acceptEncodingsStr, ",")
		if len(acceptEncodings) == 0 {
			return
		}
		var responseWriter *ResponseWriter
		switch strings.TrimSpace(acceptEncodings[0]) {
		case "gzip":
			c.Header("Content-Encoding", "gzip")
			responseWriter = NewGzip(c.Writer, gzip.BestSpeed)
		case "deflate":
			c.Header("Content-Encoding", "deflate")
			responseWriter = NewDeflate(c.Writer, flate.BestSpeed)
		case "br":
			c.Header("Content-Encoding", "br")
			responseWriter = NewBrotli(c.Writer, brotli.BestSpeed)
		default:
			return
		}
		c.Writer = responseWriter
		c.Header("Vary", "Accept-Encoding")

		defer c.Header("Content-Length", strconv.Itoa(c.Writer.Size()))
		defer responseWriter.Close()
		c.Next()
	}
}
