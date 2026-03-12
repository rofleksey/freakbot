package telegram

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Service) handleUpdates(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message != nil {
		s.handleMessage(ctx, update.Message)
	}
}

func (s *Service) handleMessage(ctx context.Context, msg *models.Message) {
	if needReply(msg.Text) ||
		msg.ReplyToMessage != nil && msg.ReplyToMessage.From != nil && msg.ReplyToMessage.From.Username == botUsername ||
		len(msg.NewChatMembers) > 0 ||
		msg.LeftChatMember != nil {

		if msg.Text == "" {
			msg.Text = "Спасибо за травлю в интернете!"
		}

		msg.Text = strings.ReplaceAll(msg.Text, "@"+botUsername, "")

		replyText, err := s.chatbotSvc.GenerateReply(ctx, msg.Text)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to generate reply", "error", err)
			return
		}

		s.tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   replyText,
			ReplyParameters: &models.ReplyParameters{
				MessageID: msg.ID,
			},
		})

		slog.InfoContext(ctx, "Freak reply", "text", replyText)
		return
	}
}
