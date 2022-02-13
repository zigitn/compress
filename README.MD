# compress

[![Go Report Card](https://goreportcard.com/badge/github.com/zigitn/compress)](https://goreportcard.com/report/github.com/zigitn/compress)

an easy use compress middleware for gin.

# usage
Use `compress.New()` to create a compress middleware, it's accept an option, you can use `compress.UseAllBestSpeed()` to use all method and set them for best speed, or use `compress.UseAllBestCompression()`
to use all method but set them for best compression. if you want highly customizable, you can define your own `Option`.  
By default, the methods order is `br`,`gzip`,`deflate`, if you like to change the order, you can define `Option` and change `Option.EnableMethods`.

# example
```go
    package main

    import (
        "github.com/gin-gonic/gin"
        "github.com/zigitn/compress"
    )

    func main() {
        router := gin.Default()
        router.Use(compress.New(compress.UseAllBestSpeed()))
        router.GET("/", func(c *gin.Context) {
            c.String(200, "Hello World")
        })
        router.Run()
    }
```

# thanks
[gin-contrib/gzip](https://github.com/gin-contrib/gzip)
[anargu/gin-brotli](https://github.com/anargu/gin-brotli)