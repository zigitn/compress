package compress

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

const (
	UseBrotli  = "br"
	UseGzip    = "gzip"
	UseDeflate = "deflate"
)

// Option can enable and config each compress method
// for simple needs please use UseAllBestSpeed or UseAllBestBestCompression
type Option struct {
	EnableMethods []string // enabled compress methods, will be matched in order
	GzipLevel     int
	BrotliOption  brotli.WriterOptions
	DeflateOption DeflateOption
}

type DeflateOption struct {
	Level int
	Dict  []byte // leave empty for default
}

func UseAllBestSpeed() Option {
	return Option{
		EnableMethods: []string{UseBrotli, UseGzip, UseDeflate},
		GzipLevel:     gzip.BestSpeed,
		BrotliOption:  brotli.WriterOptions{Quality: brotli.BestSpeed},
		DeflateOption: DeflateOption{Level: flate.BestSpeed},
	}
}

func UseAllBestBestCompression() Option {
	return Option{
		EnableMethods: []string{UseBrotli, UseGzip, UseDeflate},
		GzipLevel:     gzip.BestCompression,
		BrotliOption:  brotli.WriterOptions{Quality: brotli.BestCompression},
		DeflateOption: DeflateOption{Level: flate.BestCompression},
	}
}

// New a compress middleware, will use Option to config
func New(option Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncodingsStr := c.GetHeader("accept-encoding")
		if len(acceptEncodingsStr) == 0 {
			return
		}
		var responseWriter *ResponseWriter
		for i := range option.EnableMethods {
			if strings.Contains(acceptEncodingsStr, option.EnableMethods[i]) {
				fmt.Println(option.EnableMethods[i])
				c.Header("Content-Encoding", option.EnableMethods[i])
				methodsMap[option.EnableMethods[i]](c.Writer, option)
				break
			}
		}
		if responseWriter == nil {
			return
		}
		c.Writer = responseWriter
		c.Header("Vary", "Accept-Encoding")

		defer c.Header("Content-Length", strconv.Itoa(c.Writer.Size()))
		defer responseWriter.Close()
		c.Next()
	}
}
