package deputysheriff

import "github.com/bwmarrin/discordgo"

func init() {
	CommandsDefinitions = append(CommandsDefinitions, &discordgo.ApplicationCommand{
		Name:        "arrest-unset-channel",
		Description: "Unsets channel for annoucements.",
		// DefaultMemberPermissions: &defaultMemberPermissions,
	})
}
