package connect

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
	"gpt-bot/chat"
	"gpt-bot/client"
	"gpt-bot/json"
)

// Connect 连接到服务器
func Connect() {
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"gpt-bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:8080/", ""),
		},
	}, nil)
}

func init() {
	engine := zero.New()
	// 读取配置文件
	config, err := json.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config.json:", err)
		return
	}
	// 创建客户端
	aClient := client.CreateClient(config.Apikey, config.Proxy)
	// 创建一个哈希表
	dict := make(map[int64][]string)
	engine.OnMessage(zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		qq := ctx.Event.UserID
		if _, ok := dict[qq]; !ok {
			dict[qq] = make([]string, 2)
		}
		result, err := chat.CreateChat(aClient, config.Model, config.MaxTokens, ctx.ExtractPlainText(), dict[qq])
		if err != nil {
			fmt.Println("Error creating chat:", err)
			return
		}
		ctx.SendChain(message.Text(result))
	})
}
