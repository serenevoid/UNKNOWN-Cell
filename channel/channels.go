package channel

import (
	"fmt"
	"regexp"
	"unknown/db"
	"unknown/session"
	"unknown/utils"

	"github.com/bwmarrin/discordgo"
)

func CreateChannel() {
	messageChannel := make(chan *discordgo.MessageCreate)
	linkPattern, _ := regexp.Compile(`[a-z]+[:.].*`)
	tagPattern, _ := regexp.Compile(`.{3,32}#[0-9]{4}`)
	s := session.GetSession()
	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			if db.ViewConnection(m.ChannelID) != "" {
				if !linkPattern.MatchString(m.Content) {
					if !tagPattern.MatchString(m.Content) {
						messageChannel <- m
					}
				}
			}
		}
	})
	go func() {
		for {
			m := <-messageChannel
			// TODO: For Group Chats add number to stranger
			_, err := s.ChannelMessageSend(db.ViewConnection(m.ChannelID), fmt.Sprintf("Stranger: %s", m.Content))
			if err != nil {
				fmt.Println("error sending message: ", err)
			}
		}
	}()
}

func CreateChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	if pair := utils.GetPair(); pair != "" {
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
