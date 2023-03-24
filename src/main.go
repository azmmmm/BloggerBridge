package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"izumi.pro/wrapper/src/pkg/config"
	"izumi.pro/wrapper/src/pkg/proxy"
	"izumi.pro/wrapper/src/pkg/resource"
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
		c.String(200, "hello world")
	})

	route.GET("/proxy/*path",
		cache.CacheByRequestPath(redisStore, 30*24*time.Hour),
		proxyHandler,
	)

	route.Run(":" + strconv.Itoa(conf["port"].(int)))
}

func proxyHandler(context *gin.Context) {

	resUrl := context.Param("path")[1:] // resUrl == base64( https://www.example.com/somthine/??? )
	urlBytes, _ := base64.StdEncoding.DecodeString(resUrl)
	resUrl = string(urlBytes)

	res, err := resource.FetchByProxy(proxy.Get(), resUrl)
	if err != nil {
		log.Println("Error:", err)
	}
	sendFile(res, context)

}

// send file *os.File -> context *gin.Context
func sendFile(res *http.Response, c *gin.Context) {

	extraHeaders := map[string]string{
		//"Content-Disposition": `inline;
		//filename=` + file.Name(),
	}
	c.DataFromReader(http.StatusOK, res.ContentLength, res.Header.Get("content-type"), res.Body, extraHeaders)
}
