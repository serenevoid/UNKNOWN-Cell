package main

import (
    "log"
	"github.com/bwmarrin/discordgo"
)

var (
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with pong to confirm connectivity.",
		},
		{
			Name:        "chat",
			Description: "Connects you to a random user or a server.",
		},
		{
			Name:        "end",
			Description: "Disconnects you from the current chat.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong",
				},
			})
		},
		"chat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if pair := getPair(); pair != "" {
				pairs[i.ChannelID] = pair
				pairs[pair] = i.ChannelID
				s.ChannelMessageSend(pair, "You are connected with another user. Say Hello!")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are connected with another user. Say Hello!",
					},
				})
			} else {
				waitingUsers = append(waitingUsers, i.ChannelID)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Checking for another user...",
					},
				})
			}
		},
		"end": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if pair := pairs[i.ChannelID]; pair != "" {
				delete(pairs, i.ChannelID)
				delete(pairs, pair)
				s.ChannelMessageSend(pair, "The other user ended the chat.")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Chat ended.",
					},
				})
			}
		},
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
