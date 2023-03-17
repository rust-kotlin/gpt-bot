# GPT-Bot
这是一个能用调用openai的api接口并且能够通过[go-cqhttp](https://github.com/Mrs4s/go-cqhttp)向QQ发送消息的机器人
# 使用方法
## 1. 普通用户或非Go程序开发者
1. 下载[go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
2. 在cmd中启动`go-cqhttp`，会提示创建配置文件，选择正向WebSocket的选项，GPT-Bot默认与`go-cqhttp`通过8080端口进行WebSocket连接，请确保8080端口未被占用
3. 修改`go-cqhttp`的配置文件cofig.yml，修改其中的QQ号即可，下次启动`go-cqhttp`时可能需要扫码登录，如果登录失败，请修改device.json中的`protocol`配置项为2，具体参见该[issue](https://github.com/Mrs4s/go-cqhttp/issues/1942)
4. 在本项目[Release](https://github.com/rust-kotlin/gpt-bot/releases)页面下载最新版程序，直接打开运行，提示创建config文件，修改config.json中的配置项后再打开该程序即可
5. **config.json**文件结构介绍
- api_key: openai的api_key，可以在[openai](https://platform.openai.com/)的dashboard中找到
- model: openai的model，保持默认即可
- proxy: 本地的代理端口，通过http协议与本地代理连接，如果不通过此协议设置代理请改成空字符串
- max_tokens: 最大生成token数，越大生成的文本越长，该最大tokens限制了2个GPT的之前回答和1个最新的提问的总token数
- super_users: 超级管理员列表，只有超级管理员才可以修改GPT的角色，按以下格式填写`[12345678, 123456789]`
- weather_api: 天气api，可以在[和风天气](https://id.qweather.com/)中找到，免费api一天上限1000次查询，如不需要请留空
- base_url: 高级选择，通过cloudflare间接访问api
- temperature: 机器人的创造性，值越高越有创造性
## 2. Go程序开发者
```cmd
go install https://github.com/rust-kotlin/gpt-bot
```
如果想修改源码请获取本项目，源码
```cmd
go get https://github.com/rust-kotlin/gpt-bot
# 或者
git clone https://github.com/rust-kotlin/gpt-bot
```
# 为本项目贡献或者参考本项目思路
1. connect.go中为核心代码，调用了`ZeroBot`库与`go-cqhttp`通信，通过init()函数实现插件式，init()函数仅运行一次
2. chat.go中为与openai通信的代码，调用`go-openai`库，机器人的角色预设保存在此代码中
3. jsonconfig.go为负责读取config.json文件的代码，config.json文件中的配置项会被保存在`Config`结构体中
4. client.go负责创建openai的本地客户端，并配置使用系统代理
5. plugins下的weather.go为获取天气的插件，调用了和风天气的api
