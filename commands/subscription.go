package commands

import (
	"strconv"
	"time"
	"unknown/db"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands = append(commands,
		&discordgo.ApplicationCommand{
			Name:        "enlist",
			Description: "Enlists channel to recieve calls.",
		},
		&discordgo.ApplicationCommand{
			Name:        "delist",
			Description: "Delist channel from receiving calls.",
		},
	)
	commandHandlers["enlist"] = subscribeChannel
	commandHandlers["delist"] = unsubscribeChannel
}

func subscribeChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if db.IsKeyPresentInBucket("Channels", i.ChannelID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is already registered.",
			},
		})
		return
	}
	if i.GuildID != "" {
		if db.IsKeyPresentInBucket("Guilds", i.GuildID) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Another channel in this server has been registered earlier. Switching subscription to current channel.",
				},
			})
			return
		}
		db.InsertDataToBucket("Guilds", i.GuildID, []byte(i.ChannelID))
	}
	db.InsertDataToBucket("Channels", i.ChannelID, []byte(strconv.Itoa(time.Now().YearDay())))
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been enlisted. You will start recieving incoming calls.",
		},
	})
}

func unsubscribeChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !db.IsKeyPresentInBucket("Channels", i.ChannelID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not subsribed to unsubscribe.",
			},
		})
	}
	if i.GuildID != "" {
		db.DeleteDataFromBucket("Guilds", i.GuildID)
	}
	db.DeleteDataFromBucket("Channels", i.ChannelID)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been unsubscribed and is no longer a premium user. Incoming calls will be blocked.",
		},
	})
}
