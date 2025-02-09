package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func readyHandler(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("ready")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

	mem, err := s.GuildMembers(r.Guilds[0].ID, "", 1000)
	if err != nil {
		log.Println(err)
	}
	usernames := []string{}

	for _, m := range mem {
		usernames = append(usernames, m.User.Username)
	}
	log.Printf("Guild members: %v\n", usernames)

	err = s.RequestGuildMembers(r.Guilds[0].ID, "", 0, "", true) // Doesn't work
	if err != nil {
		log.Println(err)
	}
}
