package session

import (
  "log"
  "os"
  "os/signal"

  "github.com/bwmarrin/discordgo"
  "github.com/joho/godotenv"
)

var (
  s                *discordgo.Session
  isSessionCreated = false
)

func GetSession() *discordgo.Session {
  if !isSessionCreated {
    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }
    token, exists := os.LookupEnv("DBT")
    if !exists {
      log.Fatal("No Discord Bot Tokens found")
    }
    s, err = discordgo.New("Bot " + token)
    if err != nil {
      log.Fatal("Error creating discord session: ", err)
    } else {
      isSessionCreated = true
    }
  }
  return s
}

func CreateSocketConnection() {
  s.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
    log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
    })
  err := s.Open()
  if err != nil {
    log.Fatalf("Cannot open the session: %v", err)
  }
}

func SetupInterrupt() {
  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt)
  log.Println("Press Ctrl+C to exit")
  <-stop
}
