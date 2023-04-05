package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	cachecontrol "github.com/joeig/gin-cachecontrol/v2"
	"izumi.pro/wrapper/src/pkg/config"
	"izumi.pro/wrapper/src/pkg/proxy"
)

var TEST_CN_FLAG = true

func main() {
	conf := config.Get()
	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    conf["redis"].(string),
	}))

	route := gin.Default()

	route.GET("/", func(c *gin.Context) {
		c.String(200, "hello world 2023/3/31")
	})

	route.Use(gzip.Gzip(gzip.DefaultCompression))

	const Year = 365 * 24 * time.Hour
	route.Use(cachecontrol.New(cachecontrol.Config{
		MustRevalidate:       false,
		NoCache:              false,
		NoStore:              false,
		NoTransform:          false,
		Public:               true,
		Private:              false,
		ProxyRevalidate:      true,
		MaxAge:               cachecontrol.Duration(Year),
		SMaxAge:              nil,
		Immutable:            true,
		StaleWhileRevalidate: cachecontrol.Duration(2 * time.Hour),
		StaleIfError:         cachecontrol.Duration(2 * time.Hour),
	}))

	route.Use(cache.CacheByRequestPath(redisStore, 30*24*time.Hour))

	route.GET("/proxy/*path",

		proxyHandler,
	)

	route.Run(":" + strconv.Itoa(conf["port"].(int)))
}

func proxyHandler(context *gin.Context) {

	resUrl := context.Param("path")[1:] // resUrl == base64( https://www.example.com/somthine/??? )
	urlBytes, _ := base64.StdEncoding.DecodeString(resUrl)
	resUrl = string(urlBytes)

	res, err := proxy.FetchByProxy(resUrl)
	if err != nil {
		log.Println("Error:", err)
	}
	sendFile(res, context)

}

func sendFile(res *http.Response, c *gin.Context) {
	extraHeaders := make(map[string]string)
	c.DataFromReader(http.StatusOK, res.ContentLength, res.Header.Get("content-type"), res.Body, extraHeaders)
}
