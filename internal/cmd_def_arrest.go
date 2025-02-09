package deputysheriff

import "github.com/bwmarrin/discordgo"

func init() {
	CommandsDefinitions = append(CommandsDefinitions, &discordgo.ApplicationCommand{
		Name:        "arrest",
		Description: "Put someone to jail for some time",
		// DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Who are we arresting bestie?",
				Required:    true,
			},
		},
	})
}
