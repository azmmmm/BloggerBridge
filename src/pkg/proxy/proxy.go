package proxy

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"sync"

	"izumi.pro/wrapper/src/pkg/config"
)

var once sync.Once

var singleInstance *http.Client

func Get() *http.Client {
	once.Do(register)
	return singleInstance
}

func register() {
	conf := config.Get()

	proxyUrl, err := url.Parse(conf["proxy"].(string))
	if err != nil {
		log.Fatalf("url.Parse: %v", err)
	}
	proxyClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   10 * time.Second}

	singleInstance = proxyClient
}
