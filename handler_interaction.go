package main

import (
	ds "deputysheriff/internal"

	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand: // Registers slash commands
		if h, ok := ds.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
			defer h(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := ds.ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
			defer h(s, i)
		}
	case discordgo.InteractionModalSubmit:
		if h, ok := ds.ComponentsHandlers[i.ModalSubmitData().CustomID]; ok {
			defer h(s, i)
		}
	}
}
