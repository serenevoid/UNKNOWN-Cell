package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	waitingUsers = make([]string, 0)
	pairs        = make(map[string]string)
)

func getPair() string {
	if len(waitingUsers) < 1 {
		return ""
	}
	pair := waitingUsers[0]
	waitingUsers = waitingUsers[1:]
	return pair
}

func createChannel() {
	messageChannel := make(chan *discordgo.MessageCreate)
	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if !m.Author.Bot {
			if pairs[m.ChannelID] != "" {
				messageChannel <- m
			}
		}
	})

	go func() {
		for {
			m := <-messageChannel
            _, err := s.ChannelMessageSend(pairs[m.ChannelID], fmt.Sprintf("Stranger: %s", m.Content))
			if err != nil {
				fmt.Println("error sending message: ", err)
				return
			}
		}
	}()

}
