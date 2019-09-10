package spider

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// http 连接池管理

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		//Proxy: http.ProxyFromEnvironment, // 代理地址
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		MaxIdleConns:        200,              // 最大连接量
		MaxIdleConnsPerHost: 200,              // ip量
		IdleConnTimeout:     90 * time.Second, // 空闲time释放连接
	},
}

// 获取网页html内容
func GetHtmlBody(url string) io.ReadCloser {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Set("referer", "http://www.baidu.com/")
	req.Close = true
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res.Body
}
