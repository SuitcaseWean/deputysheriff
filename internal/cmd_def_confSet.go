package deputysheriff

import "github.com/bwmarrin/discordgo"

func init() {
	CommandsDefinitions = append(CommandsDefinitions, &discordgo.ApplicationCommand{
		Name:        "arrest-config-set",
		Description: "Arrest settings",
		// DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "annoucement-channel",
				Description: "A channel for arrest/release announcements.",
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
				Required: false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "arrest-role",
				Description: "A role that will be added to the arrested person.",
				Required:    false,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "min-time",

				Description: "Minimum time for imprisonment (Format [min]m[sec]s).",
				Required:    false,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "max-time",

				Description: "Maximum time for imprisonment (Format [min]m[sec]s).",
				Required:    false,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "default-time",

				Description: "Default time for imprisonment (Format [min]m[sec]s).",
				Required:    false,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "embed-color",

				Description: "HEX value for embed's color stripe (defaults to black).",
				Required:    false,
			},
		},
	})
}
