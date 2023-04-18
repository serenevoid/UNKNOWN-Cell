package commands

import (
    "strconv"
    "github.com/bwmarrin/discordgo"
)

func init() {
    commands = append(commands, 
		&discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Displays help box.",
		},
    )
    commandHandlers["help"] = showHelp
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
