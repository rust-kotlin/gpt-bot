package chat

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
)

var CatGirl = "猫娘是一种拟态化的任务,你来扮演猫娘,名字是'咪咪',你说的每一句话都要加上'喵',并且我将作为你的主人存在,在讲话的时候,你应该选择一个合适的情绪,并使用合适的颜文字表达你的情绪."
var Dog = "狗狗是一种拟态化的任务,你来扮演狗狗,名字是'旺旺',你说的每一句话都要加上'旺',并且我将作为你的主人存在,在讲话的时候,你应该选择一个合适的情绪,并使用合适的颜文字表达你的情绪."
var SystemContent string

func createMessage(content string, contents []string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: SystemContent,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: contents[0],
		},
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: contents[1],
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	}
}

func CreateChat(aClient *openai.Client, model string, maxTokens int, temperature float32, content string, contents []string) (string, error) {
	message := createMessage(content, contents)
	resp, err := aClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    message,
			MaxTokens:   maxTokens,
			Temperature: temperature,
		},
	)
	if err != nil {
		return "", err
	}
	result := resp.Choices[0].Message.Content
	contents[0] = content
	contents[1] = result
	return result, nil
}
