package commands

import (
	"log"
	"unknown/session"

	"github.com/bwmarrin/discordgo"
)

var (
	registeredCommands = []*discordgo.ApplicationCommand{}
	commands           = []*discordgo.ApplicationCommand{}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func AddCommandHandlers() {
    s := session.GetSession()
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func AddCommands() {
	log.Println("Adding commands...")
    s := session.GetSession()
	for _, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
        registeredCommands = append(registeredCommands, cmd)
	}
}

func RemoveCommands() {
	log.Println("Removing commands...")
    s := session.GetSession()
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Gracefully shutting down.")
}
