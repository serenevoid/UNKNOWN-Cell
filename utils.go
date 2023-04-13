package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func getUserID(i *discordgo.InteractionCreate) string {
	if i.GuildID != "" {
		return i.Member.User.ID
	} else {
		return i.User.ID
	}
}

func getPair() string {
	if len(waitingChannels) < 1 {
		return ""
	}
	pair := waitingChannels[0]
	waitingChannels = waitingChannels[1:]
	return pair
}

func setPair(user1 string, user2 string, i *discordgo.InteractionCreate) {
	pairedChannels[user1] = user2
	pairedChannels[user2] = user1
	s.ChannelMessageSend(user2, "You are connected with another user. Say Hello!")
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You are connected with another user. Say Hello!",
		},
	})
}

func unsetPair(user1 string, user2 string, s *discordgo.Session, i *discordgo.InteractionCreate) {
		delete(pairedChannels, user1)
		delete(pairedChannels, user2)
		delete(channelUserMap, user1)
		delete(channelUserMap, user2)
		s.ChannelMessageSend(user2, "The other user ended the chat.")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Chat ended.",
			},
		})
}

func addToWaitList(userID string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	waitingChannels = append(waitingChannels, i.ChannelID)
	channelUserMap[i.ChannelID] = userID
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Checking for another user...",
		},
	})
}

func isBanned(userID string, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if bannedUsers[userID] != 0 {
		if (time.Now().YearDay() - bannedUsers[userID]) > 2 {
			delete(bannedUsers, userID)
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

func isWaiting(channelID string) int {
	for index, v := range waitingChannels {
		if v == channelID {
			return index
		}
	}
	return -1
}
