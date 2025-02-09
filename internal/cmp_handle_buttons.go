package deputysheriff

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	ComponentsHandlers["unset-channel-yes"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		config.annoucementsChannel = nil
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "As you wish!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Println(err)
			sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		}
	}
	ComponentsHandlers["unset-channel-no"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Okie Dokie, the channel stays.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Println(err)
			sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		}
	}
}
