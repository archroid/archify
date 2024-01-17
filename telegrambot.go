package main

import (
	"net/http"
	"net/url"
	"os"

	utils "archroid/archify/utils"
	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func telegramBot() error {

	//setup proxy
	proxyURL, err := url.Parse("http://127.0.0.1:2081") 
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

		if update.Message.Text == "/shutdown" || update.Message.Text == "/off" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Shutting down triggered")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

			err := utils.Shutdown()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Shut down Error"+err.Error())
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

		}

		if update.Message.Text == "/reboot" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Rebooting triggered")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

			err := utils.Reboot()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reboot Error"+err.Error())
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
			
		}

		if update.Message.Text == "/sleep" || update.Message.Text == "/suspend" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Suspention triggered")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

			err := utils.Sleep()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sleep Error"+err.Error())
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
			
		}

	}
	return nil
}
