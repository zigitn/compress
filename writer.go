package compress

import (
	"io"

	"github.com/gin-gonic/gin"
)

// ResponseWriter is a generic response writer
type ResponseWriter struct {
	gin.ResponseWriter
	ComPressWriter io.WriteCloser
}

// Close closes the writer
func (r *ResponseWriter) Close() error {
	return r.ComPressWriter.Close()
}

// Write data to response
func (r *ResponseWriter) Write(data []byte) (int, error) {
	r.Header().Del("Content-Length")
	return r.ComPressWriter.Write(data)
}

// WriteHeader writes the header
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Header().Del("Content-Length")
	r.ResponseWriter.WriteHeader(statusCode)
}

// WriteString writes string data to response
func (r *ResponseWriter) WriteString(s string) (int, error) {
	r.Header().Del("Content-Length")
	return r.ComPressWriter.Write([]byte(s))
}
