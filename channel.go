package main

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	waitingChannels = make([]string, 0)
	pairedChannels  = make(map[string]string)
	channelUserMap  = make(map[string]string)
	reportedUsers   = make(map[string]int)
	bannedUsers     = make(map[string]int)
)

func getPair() string {
	if len(waitingChannels) < 1 {
		return ""
	}
	pair := waitingChannels[0]
	waitingChannels = waitingChannels[1:]
	return pair
}

func createChannel() {
	messageChannel := make(chan *discordgo.MessageCreate)
	linkPattern, _ := regexp.Compile(`[a-z]+[:.].*`)
	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			if pairedChannels[m.ChannelID] != "" {
				if !linkPattern.MatchString(m.Content) {
					messageChannel <- m
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
