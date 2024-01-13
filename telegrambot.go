package main

import (
	"net/http"
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func telegramBot() error {

	//setup proxy

	proxyURL, err := url.Parse("http://127.0.0.1:2081") // replace with your proxy address
	if err != nil {
		return err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("TELEGRAM_BOT_TOKEN"), tgbotapi.APIEndpoint, client)
	if err != nil {
		// log.Printf("Failed to create bot: %v", err)
		return err
	}
	// bot.Debug = true

	log.Info("Telegram bot is now running")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Text == "/start" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, welcome to talkbot")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
	return nil
}
