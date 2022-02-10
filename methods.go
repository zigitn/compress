package compress

import (
	"compress/flate"
	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

func NewGzip(writer gin.ResponseWriter, level int) *ResponseWriter {
	gzipWriter, err := gzip.NewWriterLevel(writer, level)
	if err != nil {
		panic(err)
	}
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: gzipWriter,
	}
}

func NewDeflate(writer gin.ResponseWriter, level int) *ResponseWriter {
	deflateWriter, err := flate.NewWriter(writer, level)
	if err != nil {
		panic(err)
	}
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: deflateWriter,
	}
}

func NewBrotli(writer gin.ResponseWriter, level int) *ResponseWriter {
	brWriter := brotli.NewWriterLevel(writer, level)
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: brWriter,
	}
}
