package main

import (
	"unknown/channel"
	"unknown/commands"
	"unknown/db"
	"unknown/session"
)

func main() {
	session.GetSession()
	commands.AddCommandHandlers()
	session.CreateSocketConnection()
	commands.AddCommands()
	defer session.GetSession().Close()
    defer db.CloseDB()
	channel.CreateChannel()
	session.SetupInterrupt()
	commands.RemoveCommands()
}
