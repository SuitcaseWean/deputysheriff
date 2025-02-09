package deputysheriff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandsHandlers["arrest-config-set"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options

		// Or convert the slice into a map
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		margs := make([]interface{}, 0, len(options))
		msgformat := "Values set:\n"

		if opt, ok := optionMap["annoucement-channel"]; ok {
			config.annoucementsChannel = opt.ChannelValue(nil)
			margs = append(margs, opt.ChannelValue(nil).ID)
			msgformat += "> annoucement-channel: <#%s>\n"
		}
		if opt, ok := optionMap["arrest-role"]; ok {
			config.arrestRole = opt.RoleValue(nil, i.GuildID)
			margs = append(margs, opt.RoleValue(nil, i.GuildID).ID)
			msgformat += "> arrest-role: <@&%s>\n"
		}
		if opt, ok := optionMap["min-time"]; ok {
			config.minTime = opt.StringValue()
			margs = append(margs, opt.StringValue())
			msgformat += "> min-time: %s\n"
		}
		if opt, ok := optionMap["max-time"]; ok {
			config.maxTime = opt.StringValue()
			margs = append(margs, opt.StringValue())
			msgformat += "> max-time: %s\n"
		}
		if opt, ok := optionMap["default-time"]; ok {
			config.defaultTime = opt.StringValue()
			margs = append(margs, opt.StringValue())
			msgformat += "> default-time: %s\n"
		}
		if opt, ok := optionMap["embed-color"]; ok {
			err := config.embedColor.colorHexToDecimal(opt.StringValue())
			if err != nil {
				fmt.Println(err)
				sendErrorResponse(s, i, errMsg(ERR_COLOR_INVALID_FORMAT, nil))

				margs = append(margs, config.embedColor.hexValue)
				msgformat += "> embed-color: `%s`\n"
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Flags: discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				})
				return
			} else {
				margs = append(margs, config.embedColor.hexValue)
				msgformat += "> embed-color: `%s`\n"
			}
		}

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
