package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client         *openai.Client
	chatModel      string
	embeddingModel openai.EmbeddingModel
}

func New(apiKey, baseURL, chatModel, embeddingModel string) *Client {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &Client{
		client:         openai.NewClientWithConfig(cfg),
		chatModel:      chatModel,
		embeddingModel: openai.EmbeddingModel(embeddingModel),
	}
}

func (c *Client) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	var all [][]float32

	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
		Input: texts,
		Model: c.embeddingModel,
	})
	if err != nil {
		return nil, fmt.Errorf("create embeddings: %w", err)
	}

	for _, e := range resp.Data {
		all = append(all, e.Embedding)
	}

	return all, nil
}

func (c *Client) ChatCompletion(ctx context.Context, system, user string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: c.chatModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: user},
		},
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("chat completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("chat completion: empty choices")
	}
	return resp.Choices[0].Message.Content, nil
}
