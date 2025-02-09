package main

import (
	"fmt"
	"log"

	ds "deputysheriff/internal"

	"github.com/bwmarrin/discordgo"
)

func commandsAdd(s *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(ds.CommandsDefinitions))
	for i, v := range ds.CommandsDefinitions {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%v' command: %v", v.Name, err)

		}
		registeredCommands[i] = cmd
	}

	return registeredCommands, nil
}
