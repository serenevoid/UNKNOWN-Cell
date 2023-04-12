package main

import (
	"log"
	"time"

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
			if (time.Now().YearDay() - bannedUsers[i.User.ID]) > 2 {
				delete(bannedUsers, i.User.ID)
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please wait till the soft ban lifts. Try to be kind from next time if you did something wrong.",
					},
				})
				return
			}
			for _, v := range waitingChannels {
				if v == i.ChannelID {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You are already in the waiting list.",
						},
					})
					return
				}
			}
			if pair := getPair(); pair != "" {
				pairedChannels[i.ChannelID] = pair
				pairedChannels[pair] = i.ChannelID
				s.ChannelMessageSend(pair, "You are connected with another user. Say Hello!")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You are connected with another user. Say Hello!",
					},
				})
			} else {
				waitingChannels = append(waitingChannels, i.ChannelID)
				channelUserMap[i.ChannelID] = i.User.ID
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Checking for another user...",
					},
				})
			}
		},
		"end": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for index, v := range waitingChannels {
				if v == i.ChannelID {
					waitingChannels = append(waitingChannels[:index], waitingChannels[index+1:]...)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Chat ended.",
						},
					})
				}
			}
			if pair := pairedChannels[i.ChannelID]; pair != "" {
				delete(pairedChannels, i.ChannelID)
				delete(pairedChannels, pair)
				delete(channelUserMap, i.ChannelID)
				s.ChannelMessageSend(pair, "The other user ended the chat.")
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Chat ended.",
					},
				})
			}
		},
		"report": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			pair := pairedChannels[i.ChannelID]
			user := channelUserMap[pair]
			reportedUsers[user] += 1
			if reportedUsers[user] > 4 {
				bannedUsers[user] = time.Now().YearDay()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The user has been reported.",
				},
			})
		},
		"reveal": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "pong",
				},
			})
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
