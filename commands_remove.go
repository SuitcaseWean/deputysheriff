package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func commandsRemove(s *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand) error {
	log.Println("Removing commands...")
	// I am not sure what is happening here or if it works
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			return fmt.Errorf("cannot delete '%v' command: %v", v.Name, err)
		}
	}
	return nil
}
