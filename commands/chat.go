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
			Description: "Connects you to a random user or a server.",
		},
		&discordgo.ApplicationCommand{
			Name:        "end",
			Description: "Disconnects you from the current chat.",
		},
		&discordgo.ApplicationCommand{
			Name:        "report",
			Description: "Reports the stranger and disconnects chat.",
		},
		&discordgo.ApplicationCommand{
			Name:        "reveal",
			Description: "Reveals the stranger's tag to you so that you can connect on discord.",
		},
	)
	commandHandlers["help"] = showHelp
	commandHandlers["chat"] = CreateChat
	commandHandlers["end"] = EndChat
	commandHandlers["report"] = ReportUser
	commandHandlers["reveal"] = RevealUser
}

func CreateChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
