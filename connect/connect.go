package connect

import (
	"fmt"
	"gpt-bot/chat"
	"gpt-bot/client"
	"gpt-bot/jsonconfig"
	"gpt-bot/plugins/weather"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var superUsers []int64

// Connect 连接到服务器
func Connect() {
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"gpt-bot"},
		CommandPrefix: "#",
		SuperUsers:    superUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:8080/", ""),
		},
	}, nil)
}

func init() {
	engine := zero.New()
	// 读取配置文件
	config, err := jsonconfig.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config.json:", err)
		return
	}
	superUsers = config.SuperUsers
	// 创建客户端
	aClient := client.CreateClient(config.Apikey, config.Proxy)
	chat.SystemContent = chat.Dog
	// 创建一个哈希表
	dict := make(map[int64][]string)
	// 帮助
	engine.OnCommand("help").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.Send(message.Text(`gpt-bot v1.6.0，一个由TomZz开发的人工智能机器人，已开源在https://github.com/rust-kotlin/gpt-bot
#help 获取帮助信息
#time 获取当前时间
#weather 获取当前天气，命令后跟着城市名称，例如：#weather北京
#role 修改人工智能角色，目前支持的角色有：狗狗，猫娘，原版，默认为狗狗，注意，仅允许超级用户修改人工智能角色`))
	})
	// 修改角色
	engine.OnCommand("role", zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		msg := ctx.MessageString()
		if strings.TrimSpace(msg) == "#role" {
			ctx.Send(message.Text("请输入角色名称，例如：#role猫娘"))
			return
		}
		role := strings.TrimSpace(strings.Replace(msg, "#role", "", 1))
		switch role {
		case "猫娘":
			chat.SystemContent = chat.CatGirl
			ctx.Send(message.Text("已修改角色为猫娘"))
		case "狗狗":
			chat.SystemContent = chat.Dog
			ctx.Send(message.Text("已修改角色为狗狗"))
		case "原版":
			chat.SystemContent = ""
			ctx.Send(message.Text("已修改角色为原版ChatGPT"))
		default:
			ctx.Send(message.Text("指定角色不存在，未进行任何修改"))
		}

	})
	// 时间
	engine.OnCommand("time").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		currentTime := time.Now()
		ctx.Send(message.Text("当前时间：" + currentTime.Format("2006年01月02日 15:04:05 ") + getChineseWeekday(currentTime.Weekday())))
	})
	// 天气
	if config.WeatherApikey != "" {
		engine.OnCommand("weather").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			msg := ctx.MessageString()
			if strings.TrimSpace(msg) == "#weather" {
				ctx.Send(message.Text("请输入城市名称，例如：#weather北京"))
				return
			}
			location := strings.TrimSpace(strings.Replace(msg, "#weather", "", 1))
			ctx.Send(message.Text("正在获取天气信息，请稍等..."))
			aWeather, err := weather.GetWeather(location, config.WeatherApikey)
			if err != nil {
				//fmt.Println("Error getting weather:", err)
				ctx.Send(message.Text(err))
				return
			}
			ctx.Send(message.Text(aWeather))
		})
	} else {
		engine.OnCommand("weather").SetBlock(true).Handle(func(ctx *zero.Ctx) {
			ctx.Send(message.Text("未配置天气API，无法获取天气信息"))
		})
	}
	// 私发消息
	engine.OnMessage(zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		qq := ctx.Event.UserID
		// 解决与Q群管家的冲突
		if qq == 2854196310 {
			return
		}
		if _, ok := dict[qq]; !ok {
			dict[qq] = make([]string, 2)
		}
		//name := ctx.GetGroupMemberInfo(ctx.Event.GroupID, qq, false).Get("nickname").Str
		result, err := chat.CreateChat(aClient, config.Model, config.MaxTokens, ctx.ExtractPlainText(), dict[qq])
		if err != nil {
			//fmt.Println("Error creating chat:", err)
			ctx.Send(message.Text("请求失败了，可能是由于网络问题或者短时间内网络请求太多，过一会再试试吧！"))
			return
		}
		if ctx.Event.GroupID != 0 {
			ctx.SendChain(message.Text(result+"\n"), message.At(qq))
		} else {
			ctx.Send(message.Text(result))
		}
	})
}

func getChineseWeekday(weekday time.Weekday) string {
	weekdays := [...]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	return weekdays[weekday]
}
