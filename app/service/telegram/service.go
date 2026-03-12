package telegram

import (
	"context"
	"freakbot/app/config"
	"freakbot/app/service/chatbot"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/samber/do"
)

type Service struct {
	tgBot      *bot.Bot
	cfg        *config.Config
	chatbotSvc *chatbot.Service
}

func New(di *do.Injector) (*Service, error) {
	tgBot := do.MustInvoke[*bot.Bot](di)

	service := &Service{
		cfg:        do.MustInvoke[*config.Config](di),
		chatbotSvc: do.MustInvoke[*chatbot.Service](di),
		tgBot:      tgBot,
	}

	tgBot.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return true
	}, service.handleUpdates)

	return service, nil
}

func (s *Service) Init(_ context.Context) {}
