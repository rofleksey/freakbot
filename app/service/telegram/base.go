package telegram

import (
	"context"
	"freakbot/app/util"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Service) handleUpdates(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message != nil {
		s.handleMessage(ctx, update.Message)
	}
}

func (s *Service) handleMessage(ctx context.Context, msg *models.Message) {
	ctx = context.WithValue(ctx, util.UsernameContextKey, msg.From.Username)
	ctx = context.WithValue(ctx, util.UserIDContextKey, msg.From.ID)
	ctx = context.WithValue(ctx, util.ChatIDContextKey, msg.Chat.ID)

	if containsBullying(msg.Text) || len(msg.NewChatMembers) > 0 || msg.LeftChatMember != nil {
		s.tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   thePhrase,
			ReplyParameters: &models.ReplyParameters{
				MessageID: msg.ID,
			},
		})
		return
	}
}
