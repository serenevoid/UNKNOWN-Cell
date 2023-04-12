package main

import (
	"log"
	"os"
    "os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var s *discordgo.Session

func createSession() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create a new Discord session using the provided bot token
	token, exists := os.LookupEnv("DBT")
	if !exists {
		log.Println("No Discord bot tokens found")
		return
	}
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}
}

func createSocketConnection() {
	s.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

func setupInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
