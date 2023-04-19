package channel

import (
	"fmt"
	"regexp"
	"strconv"
	"unknown/db"
	"unknown/session"

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
			name := "Stranger"
			if m.GuildID != "" {
				name = name + strconv.Itoa(db.GetTempUserIndex(m.ChannelID, m.Author.ID))
			}
			_, err := s.ChannelMessageSend(db.ViewConnection(m.ChannelID),
				fmt.Sprintf("%s: %s", name, m.Content))
			if err != nil {
				fmt.Println("error sending message: ", err)
			}
		}
	}()
}
