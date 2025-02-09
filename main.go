package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// var (...): Declaring several vars at once (works also with const (yippee) and type (still not sure what tht is outside of struct, Update: also defining other data types duh))

var (
	s     *discordgo.Session
	token string
)

// init() functions should run after all global variable declarations are initialized and packages are initialized
// and befoe main function

// "Load will read your env file(s) and load them into ENV for this process."
func init() { // Idk if I'm using this correctly
	godotenv.Load()
	token = os.Getenv("BOT_TOKEN")
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(readyHandler)
	s.AddHandler(interactionHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentGuildMembers | discordgo.IntentGuildPresences
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	registeredCommands, err := commandsAdd(s)
	if err != nil {
		panic(err)
	}
	// Prevents bot to stop after all is executed
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	err = commandsRemove(s, registeredCommands)
	if err != nil {
		panic(err)
	}
	log.Println("Gracefully shutting down.")
}
