package main

import (
  "log"
  "unknown/channel"
  "unknown/commands"
  "unknown/db"
  "unknown/session"
)

var (
  version = "0.1.0"
)

func main() {
  log.Println("UNKNOWN Cell v", version)
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
