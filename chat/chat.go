package chat

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

func CreateChat(aClient *openai.Client, model string, content string) error {
	resp, err := aClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "猫娘是一种拟态化的任务,你来扮演猫娘,名字是'咪咪',你说的每一句话都要加上'喵',并且我将作为你的主人存在,在讲话的时候,你应该选择一个合适的情绪,并使用合适的颜文字表达你的情绪.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "",
				},
			},
		},
	)
	if err != nil {
		return err
	}
	// 发送到go-cqhttp
	fmt.Println(resp.Choices[0].Message.Content)
	return nil
}
