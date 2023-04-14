package main

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	commands           = []*discordgo.ApplicationCommand{
		{
			Name:        "subscribers",
			Description: "Displays the total number of people subscribed to the bot.",
		},
		{
			Name:        "chat",
			Description: "Connects you to a random user or a server.",
		},
		{
			Name:        "end",
			Description: "Disconnects you from the current chat.",
		},
		{
			Name:        "report",
			Description: "Reports the stranger and disconnects chat.",
		},
		{
			Name:        "reveal",
			Description: "Reveals the stranger's tag to you so that you can connect on discord.",
		},
		{
			Name:        "help",
			Description: "Provides the list of available slash commands and their uses.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"stats":  showStats,
		"chat":   createChat,
		"end":    endChat,
		"report": reportUser,
		"reveal": revealUser,
		"help":   showHelp,
	}
)

func addCommandHandlers() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func addCommands() {
	log.Println("Adding commands...")
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func removeCommands() {
	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Gracefully shutting down.")
}

func showStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	stats := "```Total Users: " + strconv.Itoa(len(pairedChannels)) + "```"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: stats,
		},
	})
}

func showHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := ""
	for _, command := range commands {
		message = message + "`" + command.Name + "` "
	}
    s.InteractionResponseDelete(i.Interaction)
	embed := &discordgo.MessageEmbed{
		Title:       "UNKNOWN Cell",
		Description: "Welcome UNKNOWN telecom services",
		Color:       0xc0c0c0,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/icons/1096023447605358632/d40e3f9dba42ff6810535fbe64ebc1ee.webp",
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "List of commands",
				Value:  message,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Links",
				Value:  "[Invite me](https://discord.com/oauth2/authorize?client_id=1096026189811957801&permissions=19456&scope=bot) - [Support Server](https://discord.gg/mQmKudUznv)",
				Inline: false,
			},
		},
	}
	_, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		log.Fatal("Cannot send embed")
	}
}
