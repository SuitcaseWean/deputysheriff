package deputysheriff

import "github.com/bwmarrin/discordgo"

func init() {
	CommandsDefinitions = append(CommandsDefinitions, &discordgo.ApplicationCommand{
		Name:        "arrest-config-get",
		Description: "Retrieves current config for /arrest command.",
		// DefaultMemberPermissions: &defaultMemberPermissions,
	})
}
