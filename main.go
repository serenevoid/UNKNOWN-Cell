package main

import (
	"unknown/channel"
	"unknown/commands"
	"unknown/session"
)

func main() {
	session.GetSession()
	commands.AddCommandHandlers()
	session.CreateSocketConnection()
	commands.AddCommands()
	defer session.GetSession().Close()
	channel.CreateChannel()
	session.SetupInterrupt()
	commands.RemoveCommands()
}
