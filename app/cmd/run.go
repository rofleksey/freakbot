package cmd

import (
	"context"
	"freakbot/app/config"
	"freakbot/app/service/chatbot"
	"freakbot/app/service/telegram"
	"freakbot/app/util/mylog"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var configPath string

var Run = &cobra.Command{
	Use:   "run",
	Short: "Init bot",
	Run:   runBot,
}

func init() {
	Run.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to config yaml file (required)")
}

func runBot(_ *cobra.Command, _ []string) {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di := do.New()
	do.ProvideValue(di, appCtx)

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("Failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1) //nolint:gocritic
		return
	}
	do.ProvideValue(di, cfg)

	if err = mylog.Init(cfg); err != nil {
		slog.Error("Failed to init logging",
			slog.Any("error", err),
		)
		os.Exit(1)
		return
	}
	slog.InfoContext(appCtx, "Starting service...",
		slog.Bool("telegram", true),
	)

	telegramBot, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		slog.Error("Failed to create telegram bot",
			slog.Any("error", err),
		)
		os.Exit(1)
		return
	}
	do.ProvideValue(di, telegramBot)
	go telegramBot.Start(appCtx)
	defer telegramBot.Close(appCtx)

	do.Provide(di, chatbot.New)
	do.Provide(di, telegram.New)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("Shutting down server...")

		cancel()
	}()

	go do.MustInvoke[*telegram.Service](di).Init(appCtx)

	slog.Info("Listening to incoming telegram messages")
	<-appCtx.Done()

	slog.Info("Waiting for services to finish...")
	_ = di.Shutdown()
}
