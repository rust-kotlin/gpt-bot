package client

import (
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
)

func CreateClient(api string, proxy string) *openai.Client {
	// 设置代理
	aConfig := openai.DefaultConfig(api)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		fmt.Println("Error parsing proxy url:", err)
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	aConfig.HTTPClient = &http.Client{
		Transport: transport,
	}
	// 创建客户端
	return openai.NewClientWithConfig(aConfig)
}
