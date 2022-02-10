package compress

import (
	"io"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	ComPressWriter io.WriteCloser
}

func (r *ResponseWriter) Close() error {
	return r.ComPressWriter.Close()
}

func (r *ResponseWriter) Write(data []byte) (int, error) {
	r.Header().Del("Content-Length")
	return r.ComPressWriter.Write(data)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Header().Del("Content-Length")
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWriter) WriteString(s string) (int, error) {
	r.Header().Del("Content-Length")
	return r.ComPressWriter.Write([]byte(s))
}
