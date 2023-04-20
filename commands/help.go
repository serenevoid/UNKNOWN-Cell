package commands

import (
	"strconv"
	"unknown/db"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands = append(commands,
		&discordgo.ApplicationCommand{
			Name:        "about",
			Description: "Displays info about the bot.",
		},
		&discordgo.ApplicationCommand{
			Name:        "commands",
			Description: "Displays all commands and their descriptions.",
		},
		&discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Displays how the bot works.",
		},
	)
	commandHandlers["about"] = showAbout
	commandHandlers["commands"] = showCommands
	commandHandlers["help"] = showHelp
}

func showAbout(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				Name: "List of commands",
				Value: message + "\n" +
					"--------------------------------------------------------------------",
				Inline: false,
			},
			{
				Name: "Stats",
				Value: "Total listed users: " + strconv.Itoa(db.GetKeyCount("Channels")) +
					"\nTotal Servers: " + strconv.Itoa(db.GetKeyCount("Guilds")) +
					"\nTotal Active Users: " + strconv.Itoa(db.GetConnectionCount()) + "\n" +
					"--------------------------------------------------------------------",
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

func showCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := "```"
	for _, command := range commands {
		message = message + command.Name + " - " + command.Description + "\n"
	}
	message = message + "```"
	embed := &discordgo.MessageEmbed{
		Title:       "Commands",
		Description: "UNKNOWN Cell",
		Color:       0x008080,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Commands",
				Value:  message,
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

func showHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Instructions",
		Description: "UNKNOWN Cell",
		Color:       0x008080,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "What is UNKNOWN Cell?",
				Value:  "UNKNOWN Cell is a random chat bot that connects it's users together in an anonymous way. It is a personal project whose main objective was to test how go can be used to make high performant systems.",
				Inline: false,
			},
			{
				Name:   "How do I use it?",
				Value:  "To get started with UNKNOWN Cell, you can just type the command `/chat` to initiate a random chat. You will be added to a waiting list and when someone else joins, you two will be connected through a secure and anonymous tunnel. If you want to end the chat, just type `/end`. It is that easy.",
				Inline: false,
			},
			{
				Name:   "What is the enlist feature of the bot for?",
				Value:  "The bot allows it's users to enlist to the bidirectional calling feature. If you enlist to this feature, you will recieve a notification from the bot when someone else tries to start a new chat. You can enlist with `/enlist` and to opt out, just use `/delist`.",
				Inline: false,
			},
			{
				Name:   "How is this bot different from Yggdrasil which has `--userphone`?",
				Value:  "Yes this bot does have similar functionality to `--userphone`. But since this bot is specifically for random chats, it can perform with stability and the calls will not be dropped or redirected at random moments.",
				Inline: false,
			},
			{
				Name:   "Why does the bot block me from sending URLs or my discord tag through the messages?",
				Value:  "This is a choice made by the developer to ensure safety of the users.",
				Inline: false,
			},
			{
				Name:   "What should I do if I face a problem with the bot or if I have suggestions on how to make the bot better?",
				Value:  "There is a support server for the bot which is linked in `/about`. Please consider joining the server and report the issue.",
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
