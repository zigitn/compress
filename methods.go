package compress

import (
	"compress/flate"
	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

var (
	methodsMap = map[string]func(writer gin.ResponseWriter, option Option) *ResponseWriter{
		"gzip":    newGzip,
		"deflate": newDeflate,
		"br":      newBrotli,
	}
)

// newGzip returns a new gzip response writer.
func newGzip(writer gin.ResponseWriter, option Option) *ResponseWriter {
	gzipWriter, err := gzip.NewWriterLevel(writer, option.GzipLevel)
	if err != nil {
		panic(err)
	}
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: gzipWriter,
	}
}

// newDeflate returns a new deflate response writer.
func newDeflate(writer gin.ResponseWriter, option Option) *ResponseWriter {
	var deflateWriter *flate.Writer
	var err error
	if len(option.DeflateOption.Dict) == 0 {
		deflateWriter, err = flate.NewWriter(writer, option.DeflateOption.Level)
		if err != nil {
			panic(err)
		}
	} else {
		deflateWriter, err = flate.NewWriterDict(writer, option.DeflateOption.Level, option.DeflateOption.Dict)
		if err != nil {
			panic(err)
		}
	}
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: deflateWriter,
	}
}

// newBrotli returns a new brotli response writer.
func newBrotli(writer gin.ResponseWriter, option Option) *ResponseWriter {
	brWriter := brotli.NewWriterOptions(writer, option.BrotliOption)
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: brWriter,
	}
}
