package controls

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func httpsProxyConnect(req *http.Request) (*url.URL, error) {
	proxyURL, _ := url.Parse("https://proxy.ru:8443")
	proxyURL.User = url.UserPassword("login", "pass")
	return proxyURL, nil
}

func Px() {
	transport := &http.Transport{
		Proxy: httpsProxyConnect,
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://fapi.binance.com/fapi/v1/depth?limit=5&symbol=PHAUSDT", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
