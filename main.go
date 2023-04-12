package main

func main() {
	createSession()
	addCommandHandlers()
	createSocketConnection()
	addCommands()
	defer s.Close()
	createChannel()
	setupInterrupt()
	removeCommands()
}
