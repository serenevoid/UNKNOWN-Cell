package commands

import (
	"unknown/db"
	"unknown/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands = append(commands,
		&discordgo.ApplicationCommand{
			Name:        "chat",
			Description: "Connect to a random user.",
		},
		&discordgo.ApplicationCommand{
			Name:        "end",
			Description: "Disconnect from the current chat.",
		},
		&discordgo.ApplicationCommand{
			Name:        "report",
			Description: "Report the stranger and disconnect the chat.",
		},
		&discordgo.ApplicationCommand{
			Name:        "reveal",
			Description: "Reveals your discord tag to the stranger.",
		},
	)
	commandHandlers["help"] = showHelp
	commandHandlers["chat"] = CreateChat
	commandHandlers["end"] = EndChat
	commandHandlers["report"] = ReportUser
	commandHandlers["reveal"] = RevealUser
}

func CreateChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // Verify if a connection already exists for the channel
	if db.ViewConnection(i.ChannelID) != "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cannot start a new chat till the current chat ends.",
			},
		})
		return
	}
	userID := utils.GetUserID(i)
	if db.IsBanned(userID) {
		return
	}
    // Check if the channel is in the waiting list
	if db.IsWaiting(i.ChannelID) != -1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are already in the waiting list.",
			},
		})
		return
	}
	if pair := db.PopWaitList(); pair != "" {
		utils.SetPair(i.ChannelID, pair, i)
	} else {
		db.PushWaitList(i.ChannelID)
        if db.IsKeyPresentInBucket("Channels", i.ChannelID) {
            db.GetRandomSubscribers(i.ChannelID, func(channelID string) {
                s.ChannelMessageSend(channelID, "*beep beep*\nA random user is trying to connect. To respond, type the command `/chat`.")
            })
        }
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please wait till another user hops in.",
			},
		})
	}
}

func EndChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if index := db.IsWaiting(i.ChannelID); index != -1 {
		db.RemoveWaitList(index)
		db.RemoveChannelUser(i.ChannelID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Chat ended.",
			},
		})
	}
	if pair := db.ViewConnection(i.ChannelID); pair != "" {
		utils.UnsetPair(i.ChannelID, pair, i)
	}
	if i.GuildID != "" {
		db.RemoveTempUsers(i.ChannelID)
	}
}

func RevealUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if pair := db.ViewConnection(i.ChannelID); pair != "" {
		s.ChannelMessageSend(pair, "The stranger has chosen to reveal their discord tag to you.")
		tag := utils.GetUserTag(i)
		s.ChannelMessageSend(pair, tag)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Your Discord tag has been releaved to the stranger.",
			},
		})
	}
}

func ReportUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	pair := db.ViewConnection(i.ChannelID)
	user := db.ViewChannelUser(i.ChannelID)
	db.ReportUser(user)
	if pair != "" {
		utils.UnsetPair(i.ChannelID, pair, i)
		s.ChannelMessageSend(pair, "The other user ended the chat.")
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The user has been reported.",
		},
	})
}
