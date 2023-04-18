package utils

import (
	"time"

	"unknown/db"
	"unknown/session"

	"github.com/bwmarrin/discordgo"
)

func GetUserID(i *discordgo.InteractionCreate) string {
	if i.GuildID != "" {
		return i.Member.User.ID
	} else {
		return i.User.ID
	}
}

func GetUserTag(i *discordgo.InteractionCreate) string {
	if i.GuildID != "" {
		return i.Member.User.Username+"#"+i.Member.User.Discriminator
	} else {
		return i.User.Username+"#"+i.User.Discriminator
	}
}

func GetPair() string {
	return db.PopWaitList()
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

// TODO: Check for Implementation and switch to DB
func IsBanned(banList map[string]int, userID string, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if banList[userID] != 0 {
		if (time.Now().YearDay() - banList[userID]) > 2 {
			delete(banList, userID)
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please wait till the soft ban lifts. Try to be kind from next time if you did something wrong.",
				},
			})
			return true
		}
	}
	return false
}
