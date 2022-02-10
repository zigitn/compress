package compress

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncodingsStr := c.GetHeader("accept-encoding")
		if len(acceptEncodingsStr) == 0 {
			c.Next()
		}
		acceptEncodings := strings.Split(acceptEncodingsStr, ", ")
		if len(acceptEncodings) == 0 {
			c.Next()
		}
		switch acceptEncodings[0] {
		case "gzip":
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")
			writer, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Writer = &ResponseWriter{c.Writer, writer}
			defer func() {
				_ = writer.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
			c.Next()
		case "deflate":
			c.Header("Content-Encoding", "deflate")
			c.Header("Vary", "Accept-Encoding")
			writer, err := flate.NewWriter(c.Writer, flate.BestSpeed)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Writer = &ResponseWriter{c.Writer, writer}
			defer func() {
				_ = writer.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
			c.Next()
		case "br":
			c.Header("Content-Encoding", "br")
			c.Header("Vary", "Accept-Encoding")
			writer := brotli.NewWriterLevel(c.Writer, brotli.BestSpeed)
			c.Writer = &ResponseWriter{c.Writer, writer}
			defer func() {
				_ = writer.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
			c.Next()
		case "compress":
			c.Header("Content-Encoding", "compress")
			c.Header("Vary", "Accept-Encoding")
			writer := lzw.NewWriter(c.Writer, lzw.LSB, 8)
			c.Writer = &ResponseWriter{c.Writer, writer}
			defer func() {
				_ = writer.Close()
				c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
			}()
			c.Next()
		default:
			c.Next()
		}
	}
}

type ResponseWriter struct {
	gin.ResponseWriter
	ComPressWriter io.WriteCloser
}

func (r *ResponseWriter) Write(data []byte) (int, error) {
	return r.ComPressWriter.Write(data)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Header().Del("Content-Length")
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWriter) WriteString(s string) (int, error) {
	return r.ComPressWriter.Write([]byte(s))
}
