package proxy

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"izumi.pro/wrapper/src/pkg/config"
)

var singletonProxyClient *http.Client

func init() {
	conf := config.Get()

	proxyUrl, err := url.Parse(conf["proxy"].(string))
	if err != nil {
		log.Fatalf("url.Parse: %v", err)
	}
	proxyClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		Timeout:   15 * time.Second}

	singletonProxyClient = proxyClient
}

func FetchByProxy(url string) (*http.Response, error) {
	res, err := singletonProxyClient.Get(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Printf("[%v] Proxy fetching %s", res.StatusCode, url)
	return res, nil
}
