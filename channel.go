package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	waitingChannels = make([]string, 0)
	pairedChannels  = make(map[string]string)
	channelUserMap  = make(map[string]string)
	reportedUsers   = make(map[string]int)
	bannedUsers     = make(map[string]int)
)

func createChannel() {
	messageChannel := make(chan *discordgo.MessageCreate)
	linkPattern, _ := regexp.Compile(`[a-z]+[:.].*`)
	tagPattern, _ := regexp.Compile(`.{3,32}#[0-9]{4}`)
	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			if pairedChannels[m.ChannelID] != "" {
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
			_, err := s.ChannelMessageSend(pairedChannels[m.ChannelID], fmt.Sprintf("Stranger: %s", m.Content))
			if err != nil {
				fmt.Println("error sending message: ", err)
				return
			}
		}
	}()
}

func createChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := getUserID(i)
	if isBanned(userID, s, i) {
		return
	}
	if isWaiting(i.ChannelID) != -1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are already in the waiting list.",
			},
		})
		return
	}
	if pair := getPair(); pair != "" {
		setPair(i.ChannelID, pair, i)
	} else {
		addToWaitList(userID, s, i)
	}
}

func endChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if index := isWaiting(i.ChannelID); index != -1 {
		waitingChannels = append(waitingChannels[:index], waitingChannels[index+1:]...)
		delete(channelUserMap, i.ChannelID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Chat ended.",
			},
		})
	}
	if pair := pairedChannels[i.ChannelID]; pair != "" {
		unsetPair(i.ChannelID, pair, s, i)
	}
}

func revealUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if pair := pairedChannels[i.ChannelID]; pair != "" {
		s.ChannelMessageSend(pair, "The stranger has chosen to reveal their discord tag to you.")
		if i.GuildID != "" {
			s.ChannelMessageSend(pair, i.Member.User.Username+"#"+i.Member.User.Discriminator)
		} else {
			s.ChannelMessageSend(pair, i.User.Username+"#"+i.User.Discriminator)
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Your Discord tag has been releaved to the stranger.",
			},
		})
	}
}

func reportUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	pair := pairedChannels[i.ChannelID]
	user := channelUserMap[pair]
	reportedUsers[user] += 1
	if reportedUsers[user] > 4 {
		bannedUsers[user] = time.Now().YearDay()
	}
	if pair := pairedChannels[i.ChannelID]; pair != "" {
		unsetPair(i.ChannelID, pair, s, i)
		s.ChannelMessageSend(pair, "The other user ended the chat.")
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The user has been reported.",
		},
	})
}
