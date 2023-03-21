package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"izumi.pro/wrapper/src/pkg/config"
	"izumi.pro/wrapper/src/pkg/proxy"
	"izumi.pro/wrapper/src/pkg/resource"
)

var TEST_CN_FLAG = true

func main() {
	conf := config.Get()
	log.Printf("program start!!!!!!!!!")
	route := gin.Default()
	route.GET("/", showIndex)
	route.GET("/proxy/*path", responseByCountry)

	route.Run(":" + strconv.Itoa(conf["port"].(int)))
}
func showIndex(c *gin.Context) {
	c.String(200, "hi!")
}
func responseByCountry(c *gin.Context) {
	c.RemoteIP()
	info := get_ip_info(c.RemoteIP())

	resUrl := c.Param("path")[1:] // resUrl == base64( https://www.example.com/somthine/??? )
	urlBytes, _ := base64.StdEncoding.DecodeString(resUrl)
	resUrl = string(urlBytes)

	if TEST_CN_FLAG || info["country"] == "CN" {

		log.Printf("[%v] Proxy IP source to %v", info["country"], resUrl)

		res, err := resource.FetchByProxy(proxy.Get(), resUrl)
		if err != nil {
			log.Println("Error:", err)
		}
		sendFile(res, c)

	} else {
		log.Printf("[%v] IP sourceRedirect to %v", info["country"], resUrl)
		c.Redirect(http.StatusFound, resUrl)
	}

}

// send file *os.File -> context *gin.Context
func sendFile(res *http.Response, c *gin.Context) {

	extraHeaders := map[string]string{
		//"Content-Disposition": `inline;
		//filename=` + file.Name(),
	}
	c.DataFromReader(http.StatusOK, res.ContentLength, res.Header.Get("content-type"), res.Body, extraHeaders)
}

func get_ip_info(ip string) map[string]string {

	/*	curl ipinfo.io/47.100.131.168
		{
		  "ip": "47.100.131.168",
		  "city": "Shanghai",
		  "region": "Shanghai",
		  "country": "CN",
		  "loc": "31.2222,121.4581",
		  "org": "AS37963 Hangzhou Alibaba Advertising Co.,Ltd.",
		  "timezone": "Asia/Shanghai",
		  "readme": "https://ipinfo.io/missingauth"
		}*/
	proxyClient := proxy.Get()
	resp, err := proxyClient.Get("http://ipinfo.io/" + ip)
	if err != nil {
		log.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var info map[string]string
	json.Unmarshal(body, &info)
	return info
}
