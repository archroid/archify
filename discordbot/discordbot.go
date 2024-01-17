package discordbot

import (
	"os"
	"os/signal"
	"syscall"

	utils "archroid/archify/utils"

	"github.com/bwmarrin/discordgo"
	log "github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var GuildID = "801840788022624296"

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

func RunSession() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	// Create bot session
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return err
	}
	// add intents
	dg.Identify.Intents = discordgo.MakeIntent(Intents)

	// add handlers
	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)

	// dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
	// 		h(s, i)
	// 	}
	// })

	err = dg.Open()
	if err != nil {
		return err
	}
	// adding slash commands
	// registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	// for i, v := range commands {
	// 	cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
	// 	if err != nil {
	// 		log.Errorf("Cannot create '%v' command: %v", v.Name, err)
	// 	}
	// 	registeredCommands[i] = cmd
	// }

	log.Info("Discord bot is now running")

	// Cleanly close down the Discord session.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// for _, v := range registeredCommands {
	// 	log.Info("Deleting command: ", v.Name)
	// 	err := dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, v.ID)
	// 	if err != nil {
	// 		log.Errorf("Cannot delete '%v' command: %v", v.Name, err)
	// 	}
	// }

	dg.Close()
	return nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateListeningStatus("/help")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "shutdown" || m.Content == "off" {
		s.ChannelMessageSend(m.ChannelID, "Shutting Down")

		err := utils.Shutdown()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error shutting down: "+err.Error())
			return
		}
	}

	if m.Content == "sleep" || m.Content == "suspend" {
		s.ChannelMessageSend(m.ChannelID, "Sleeping")

		err := utils.Sleep()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error sleeping: "+err.Error())
			return
		}
	}

	if m.Content == "reboot" {
		s.ChannelMessageSend(m.ChannelID, "Reeboting")

		err := utils.Reboot()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error rebooting: "+err.Error())
			return
		}
	}

}

// var (
// 	commands = []*discordgo.ApplicationCommand{
// 		{
// 			Name:        "ping",
// 			Description: "ping",
// 		},
// 	}

// 	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionResponseChannelMessageWithSource,
// 				Data: &discordgo.InteractionResponseData{
// 					Content: "pong",
// 				},
// 			})
// 		},
// 	}
// )
