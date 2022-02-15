package compress

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

const (
	UseBrotli  = "br"      // br
	UseGzip    = "gzip"    // gzip
	UseDeflate = "deflate" // deflate
)

// Option can enable and config each compress method
// for simple needs please use UseAllBestSpeed or UseAllBestBestCompression
type Option struct {
	EnableMethods    []string // enabled compress methods, will be matched in order
	GzipLevel        int
	BrotliOption     brotli.WriterOptions
	DeflateOption    DeflateOption
	HeadFilter       map[string]string
	ExtensionsFilter []string
	PathFilter       []string
	PathRegexFilter  []string
	CustomFilter     []func(c *gin.Context) bool // return true to exclude anything you want
}

// DeflateOption is used to config deflate
type DeflateOption struct {
	Level int
	Dict  []byte // leave empty for default
}

// UseAllBestSpeed will enable all compress methods with best speed
func UseAllBestSpeed() Option {
	return Option{
		EnableMethods: []string{UseBrotli, UseGzip, UseDeflate},
		GzipLevel:     gzip.BestSpeed,
		BrotliOption:  brotli.WriterOptions{Quality: brotli.BestSpeed},
		DeflateOption: DeflateOption{Level: flate.BestSpeed},
		HeadFilter: map[string]string{
			"Connection": "Upgrade",
			"Accept":     "text/event-stream",
		},
		ExtensionsFilter: []string{".png", ".gif", ".jpeg", ".jpg", ".svg", ".woff", ".woff2", ".ttf", ".eot"},
	}
}

// UseAllBestBestCompression will enable all compress methods with best compression
func UseAllBestBestCompression() Option {
	return Option{
		EnableMethods: []string{UseBrotli, UseGzip, UseDeflate},
		GzipLevel:     gzip.BestCompression,
		BrotliOption:  brotli.WriterOptions{Quality: brotli.BestCompression},
		DeflateOption: DeflateOption{Level: flate.BestCompression},
		HeadFilter: map[string]string{
			"Connection": "Upgrade",
			"Accept":     "text/event-stream",
		},
		ExtensionsFilter: []string{".png", ".gif", ".jpeg", ".jpg", ".svg", ".woff", ".woff2", ".ttf", ".eot"},
	}
}

// New a compress middleware, will use Option to config
func New(option Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check extension
		for _, extension := range option.ExtensionsFilter {
			if strings.HasSuffix(c.Request.URL.Path, extension) {
				return
			}
		}
		// check path
		for _, path := range option.PathFilter {
			if c.Request.URL.Path == path {
				return
			}
		}
		// check path regex
		for _, path := range option.PathRegexFilter {
			if match, _ := regexp.MatchString(path, c.Request.URL.Path); match {
				return
			}
		}
		// custom check func
		for _, f := range option.CustomFilter {
			if f(c) {
				return
			}
		}

		for k, v := range option.HeadFilter {
			if strings.Contains(c.GetHeader(k), v) {
				return
			}
		}
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
