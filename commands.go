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
			Name:        "register",
			Description: "Registers current channel for the Cell.",
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
			Description: "Displays help box.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "register":     registerChannel,
		"chat":        createChat,
		"end":         endChat,
		"report":      reportUser,
		"reveal":      revealUser,
		"help":        showHelp,
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
	stats := "```Total Subscribers: " + strconv.Itoa(len(pairedChannels)) + 
    "\nTotal Active Users: " + strconv.Itoa(len(pairedChannels)) + "```"
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
	embed := &discordgo.MessageEmbed{
		Title:       "UNKNOWN Cell",
		Description: "Welcome UNKNOWN telecom services.\nWe connect your calls to random users.",
		Color:       0x008080,
        Image: &discordgo.MessageEmbedImage{
            URL: "https://cdn.discordapp.com/app-icons/962387295250563092/ff587500912ce378b6672aa7a4997cd4.png",
        },
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "List of commands",
				Value:  message,
				Inline: false,
			},
			{
				Name:   "Stats",
                Value:  "Total Users: " + strconv.Itoa(0) + 
                "\nTotal Active Users: " + strconv.Itoa(0) + "\n",
				Inline: false,
			},
			{
				Name:   "Links",
				Value:  "[Invite me](https://discord.com/oauth2/authorize?client_id=1096026189811957801&permissions=19456&scope=bot) - [Support Server](https://discord.gg/mQmKudUznv)",
				Inline: false,
			},
		},
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func registerChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	stats := "```Total Subscribers: " + strconv.Itoa(len(pairedChannels)) + 
    "\nTotal Active Users: " + strconv.Itoa(len(pairedChannels)) + "```"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: stats,
		},
	})
}
