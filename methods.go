package compress

import (
	"compress/flate"
	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

var (
	methodsMap = map[string]func(writer gin.ResponseWriter, option Option) *ResponseWriter{
		"gzip":    NewGzip,
		"deflate": NewDeflate,
		"br":      NewBrotli,
	}
)

func NewGzip(writer gin.ResponseWriter, option Option) *ResponseWriter {
	gzipWriter, err := gzip.NewWriterLevel(writer, option.GzipLevel)
	if err != nil {
		panic(err)
	}
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: gzipWriter,
	}
}

func NewDeflate(writer gin.ResponseWriter, option Option) *ResponseWriter {
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

func NewBrotli(writer gin.ResponseWriter, option Option) *ResponseWriter {
	brWriter := brotli.NewWriterOptions(writer, option.BrotliOption)
	return &ResponseWriter{
		ResponseWriter: writer,
		ComPressWriter: brWriter,
	}
}
