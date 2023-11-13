package channel

import (
  "fmt"
  "regexp"
  "unknown/db"
  "unknown/session"

  "github.com/bwmarrin/discordgo"
)

var (
  UserIcons = []string{
    "dog", "cat", "mouse", "hamster", "rabbit", "fox", "bear", "panda_face",
    "polar_bear", "koala", "tiger", "lion", "cow", "pig", "frog", "monkey_face",
    "chicken", "wolf", "boar", "octopus", "squid", "gorilla", "deer", "raccoon",
    "hedgehog",
  }
)

func CreateChannel() {
  messageChannel := make(chan *discordgo.MessageCreate, 1000)
  linkPattern, _ := regexp.Compile(`[a-z]+[:.].*`)
  tagPattern, _ := regexp.Compile(`.{3,32}#[0-9]{4}`)
  s := session.GetSession()
  s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
    if !m.Author.Bot {
      if db.ViewConnection(m.ChannelID) != "" {
        if !linkPattern.MatchString(m.Content) {
          if !tagPattern.MatchString(m.Content) {
            messageChannel <- m
          } else {
            s.ChannelMessageSend(m.ChannelID, "Please do not share tags over the call. Type `/reveal` to reveal your user tag.")
          }
        } else {
          s.ChannelMessageSend(m.ChannelID, "Please do not share any links over the call.")
        }
      }
    }
  })
  go func() {
    for {
      m := <-messageChannel
      go sendMessage(s, m)
    }
  }()
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  name := ":" + UserIcons[db.GetTempUserIndex(m.ChannelID, m.Author.ID)] + ":"
  _, err := s.ChannelMessageSend(db.ViewConnection(m.ChannelID),
    fmt.Sprintf("%s **:** %s", name, m.Content))
  if err != nil {
    fmt.Println("error sending message: ", err)
  }
}
