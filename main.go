package main

import (
	"fmt"
	"gpt-bot/chat"
	"gpt-bot/client"
	"gpt-bot/json"
)

func main() {
	// 读取配置文件
	config, err := json.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config.json:", err)
		return
	}
	// 创建客户端
	aClient := client.CreateClient(config.Apikey, config.Proxy)
	// 创建聊天
	if err := chat.CreateChat(aClient, config.Model, "你好啊"); err != nil {
		fmt.Println("Error creating chat:", err)
		return
	}
}
