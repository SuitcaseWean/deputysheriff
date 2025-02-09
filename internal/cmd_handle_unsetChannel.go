package deputysheriff

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandsHandlers["arrest-unset-channel"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Are you sure you want to unset the announcement channel? If not set, announcements will show up in the channel where /arrest is called.",
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					// ActionRow is a container of all buttons within the same row.
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								// Label is what the user will see on the button.
								Label: "Yes",
								// Style provides coloring of the button. There are not so many styles tho.
								Style: discordgo.SuccessButton,
								// Disabled allows bot to disable some buttons for users.
								Disabled: false,
								// CustomID is a thing telling Discord which data to send when this button will be pressed.
								CustomID: "unset-channel-yes", // Handled in commands_handlers
							},
							discordgo.Button{
								Label:    "No",
								Style:    discordgo.DangerButton,
								Disabled: false,
								CustomID: "unset-channel-no",
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
			sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		}
	}
}
