package main

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	commands           = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with pong to confirm connectivity.",
		},
		{
			Name:        "stats",
			Description: "Displays the stats about the bot.",
		},
		{
			Name:        "chat",
			Description: "Connects you to a random user or a server.",
		},
		{
			Name:        "end",
			Description: "Disconnects you from the current chat.",
		},
		{
			Name:        "report",
			Description: "Reports the stranger and disconnects chat.",
		},
		{
			Name:        "reveal",
			Description: "Reveals the stranger's tag to you so that you can connect on discord.",
		},
		{
			Name:        "help",
			Description: "Provides the list of available slash commands and their uses.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":   ping,
		"stats":   showStats,
		"chat":   createChat,
		"end":    endChat,
		"report": reportUser,
		"reveal": revealUser,
		"help":   showHelp,
	}
)

func addCommandHandlers() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func addCommands() {
	log.Println("Adding commands...")
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func removeCommands() {
	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Gracefully shutting down.")
}

func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})
}

func showStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
    stats := "```Total Users: " + strconv.Itoa(len(pairedChannels)) + "```"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: stats,
		},
	})
}

func showHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
    message := "**List of available commands**\n```"
	for index, command := range commands {
        message = message + command.Name + " - " + command.Description
        if index != len(commands) - 1 {
            message = message + "\n"
        } else {
            message = message + "```"
        }
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
