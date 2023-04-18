package commands

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands = append(commands,
		&discordgo.ApplicationCommand{
			Name:        "subscribe",
			Description: "Subscribe to get good chance of getting connections. Will start receiving calls.",
		},
		&discordgo.ApplicationCommand{
			Name:        "unsubscribe",
			Description: "Unsubscribe from receiving calls.",
		},
	)
    commandHandlers["subscribe"] = subscribeChannel
    commandHandlers["unsubscribe"] = unsubscribeChannel
}

func subscribeChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been upgraded to a premium user. You will start recieving incoming calls.",
		},
	})
}

func unsubscribeChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have been unsubscribed and is no longer a premium user. Incoming calls will be blocked.",
		},
	})
}
