package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/charmbracelet/log"
)

func discordBot() error {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return err
	}
	const (
		Intents = discordgo.IntentsDirectMessages |
			discordgo.IntentsGuildBans |
			discordgo.IntentsGuildEmojis |
			discordgo.IntentsGuildIntegrations |
			discordgo.IntentsGuildInvites |
			discordgo.IntentsGuildMembers |
			discordgo.IntentsGuildMessageReactions |
			discordgo.IntentsGuildMessages |
			discordgo.IntentsGuildVoiceStates |
			discordgo.IntentsGuilds |
			discordgo.IntentsGuildVoiceStates
	)

	dg.Identify.Intents = discordgo.MakeIntent(Intents)

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Discord bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	return nil

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "shutdown" || m.Content == "off" {
		s.ChannelMessageSend(m.ChannelID, "Shutting Down")

		err := shutdown()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error shutting down: "+err.Error())
			return
		}

	}

	if m.Content == "sleep" || m.Content == "suspend" {
		s.ChannelMessageSend(m.ChannelID, "Sleeping")

		err := sleep()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error sleeping: "+err.Error())
			return
		}
	}

	if m.Content == "reboot" {
		s.ChannelMessageSend(m.ChannelID, "Reeboting")

		err := reboot()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error rebooting: "+err.Error())
			return
		}
	}

	

}
