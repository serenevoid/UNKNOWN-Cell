package utils

import (
  "unknown/db"
  "unknown/session"

  "github.com/bwmarrin/discordgo"
)

func GetUserID(i *discordgo.InteractionCreate) string {
  return i.Member.User.ID
}

func GetUserTag(i *discordgo.InteractionCreate) string {
  return i.Member.User.Username+"#"+i.Member.User.Discriminator
}

func SetPair(user1 string, user2 string, i *discordgo.InteractionCreate) {
  db.AddConnection(user1, user2)
  s := session.GetSession()
  s.ChannelMessageSend(user2, "You are connected with another user. Say Hello!")
  s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: "You are connected with another user. Say Hello!",
    },
    })
}

func UnsetPair(user1 string, user2 string, i *discordgo.InteractionCreate) {
  db.RemoveConnection(user1, user2)
  db.RemoveChannelUser(user1)
  db.RemoveChannelUser(user2)
  s := session.GetSession()
  s.ChannelMessageSend(user2, "The other user ended the chat.")
  s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: "Chat ended.",
    },
    })
}

func AddToWaitList(userID string, i *discordgo.InteractionCreate) {
  db.PushWaitList(i.ChannelID)
  db.AddChannelUser(i.ChannelID, userID)
  s := session.GetSession()
  s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: "Checking for another user...",
    },
    })
}
