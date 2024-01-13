package main

import (
	"os"
	"os/exec"
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
		s.ChannelMessageSend(m.ChannelID, "SHUTTING DOWN")
		// Execute the shutdown command
		cmd := exec.Command("shutdown", "-h", "now")
		err := cmd.Run()
		if err != nil {
			log.Error("Error shutting down:", err)
			// http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// If the command executed successfully, terminate the server
		log.Error("System shutdown triggered")
		os.Exit(0)

	}

	if m.Content == "sleep" || m.Content == "suspend" {
		s.ChannelMessageSend(m.ChannelID, "SLEEPING")
		// Execute the sleep command
		cmd := exec.Command("systemctl", "suspend")
		err := cmd.Run()
		if err != nil {
			log.Error("Error sleeping:", err)
			// http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Error("System suspend triggered")
	}

}
