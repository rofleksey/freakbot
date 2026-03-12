package chatbot

import (
	"context"
	"errors"
	"fmt"
	"freakbot/app/config"
	"freakbot/app/service/chatbot/knowledge"
	"freakbot/app/service/chatbot/llm"
	"strings"

	"github.com/samber/do"
)

type Service struct {
	systemPrompt string
	db           *knowledge.DB
	client       *llm.Client
	cfg          *config.Config
}

func New(di *do.Injector) (*Service, error) {
	cfg := do.MustInvoke[*config.Config](di)

	systemPrompt, db, err := knowledge.Load("data")
	if err != nil {
		return nil, fmt.Errorf("knowledge.Load: %w", err)
	}

	client := llm.New(
		cfg.OpenAI.APIKey,
		cfg.OpenAI.BaseURL,
		cfg.OpenAI.ChatModel,
		cfg.OpenAI.EmbeddingModel,
	)

	return &Service{
		systemPrompt: systemPrompt,
		db:           db,
		client:       client,
		cfg:          cfg,
	}, nil
}

func (s *Service) GenerateReply(ctx context.Context, msgText string) (string, error) {
	embeddings, err := s.client.Embed(ctx, []string{msgText})
	if err != nil {
		return "", fmt.Errorf("s.client.Embed: %w", err)
	}
	if len(embeddings) == 0 {
		return "", errors.New("no embeddings")
	}

	indices := s.db.TopKSimilar(embeddings[0], s.cfg.Retrieval.TopK)

	var b strings.Builder
	b.WriteString("Текущий вопрос пользователя:\n")
	b.WriteString(msgText)
	b.WriteString("\n\n")
	b.WriteString("Ниже список возможных ответов из истории чата. Каждый вариант — реальное сообщение участника.\n")
	b.WriteString("Формат: [CANDIDATE N] автор | дата — текст сообщения.\n\n")
	for i, idx := range indices {
		m := s.db.Messages[idx]
		fmt.Fprintf(&b, "[CANDIDATE %d] %s | %s — %s\n\n", i+1, m.From, m.Date, m.Text)
	}
	b.WriteString("Выбери строго ОДНО сообщение, которое лучше всего отвечает на вопрос и соответствует стилю чата.\n")
	b.WriteString("Верни СТРОГО только текст выбранного сообщения, без номера, без комментариев и без каких‑либо изменений или добавлений.")

	reply, err := s.client.ChatCompletion(ctx, s.systemPrompt, b.String())
	if err != nil {
		return "", fmt.Errorf("s.client.ChatCompletion: %w", err)
	}

	return strings.TrimSpace(reply), nil
}
