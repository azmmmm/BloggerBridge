package resource

import (
	"log"
	"net/http"
)

func FetchByProxy(proxyClient *http.Client, url string) (*http.Response, error) {

	res, err := proxyClient.Get(url)
	log.Printf("[%v] Proxy fetching %s", res.StatusCode, url)
	if err != nil {
		log.Print(err)
	}
	return res, err

}
