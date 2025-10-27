package telegram

import (
	"context"
	"freakbot/app/config"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/samber/do"
)

type Service struct {
	tgBot *bot.Bot
	cfg   *config.Config
}

func New(di *do.Injector) (*Service, error) {
	tgBot := do.MustInvoke[*bot.Bot](di)

	service := &Service{
		cfg:   do.MustInvoke[*config.Config](di),
		tgBot: tgBot,
	}

	tgBot.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return true
	}, service.handleUpdates)

	return service, nil
}

func (s *Service) Init(_ context.Context) {}
