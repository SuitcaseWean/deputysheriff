package deputysheriff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandsHandlers["arrest-config-get"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		margs := make([]interface{}, 0)
		msgformat := "Values set:\n"
		if config.annoucementsChannel == nil {
			msgformat += "> annoucement-channel: unset\n"
		} else {
			margs = append(margs, config.annoucementsChannel.ID)
			msgformat += "> annoucement-channel: <#%s>\n"
		}
		if config.arrestRole == nil {
			msgformat += "> arrest-role: unset\n"
		} else {
			margs = append(margs, config.arrestRole.ID)
			msgformat += "> arrest-role: <@&%s>\n"
		}
		margs = append(margs, config.minTime)
		msgformat += "> min-time: %s\n"
		margs = append(margs, config.maxTime)
		msgformat += "> max-time: %s\n"
		margs = append(margs, config.defaultTime)
		msgformat += "> default-time: %s\n"
		margs = append(margs, config.embedColor.hexValue)
		msgformat += "> embed-color: `%s`\n"

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf(
					msgformat,
					margs...,
				),
			},
		})
	}
}
