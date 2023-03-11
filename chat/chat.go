package chat

import (
	"context"
	"errors"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
)

func CreateChat(aClient *openai.Client, model string, prompt string) error {
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     model,
		MaxTokens: 100,
		Prompt:    prompt,
		Stream:    true,
	}
	stream, err := aClient.CreateCompletionStream(ctx, req)
	if err != nil {
		fmt.Println("Error creating stream:", err)
		return err
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return nil
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}

		fmt.Printf("Stream response: %v\n", response)
	}
}
