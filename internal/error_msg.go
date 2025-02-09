package deputysheriff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	ERR_DEFAULT              int = 1
	ERR_TIME_BELOW_MIN       int = 2
	ERR_TIME_ABOVE_MAX       int = 3
	ERR_TIME_INVALID_FORMAT  int = 4
	ERR_USER_ALREADY_IN_JAIL int = 5
	ERR_COLOR_INVALID_FORMAT int = 6
	ERR_SOMETHING_WENT_WRONG int = 7
)

func errMsg(errCode int, selectedUser *discordgo.User) string {
	responseMsg := ""

	switch errCode {
	case ERR_TIME_INVALID_FORMAT:
		responseMsg = fmt.Sprintln("You entered invalid tim\nThe format should look like this: [minutes]m[seconds]s.\nFor example these are valid time formats: 5m, 30s, 1m20s, 0m50s, 2m0s..")
	case ERR_TIME_BELOW_MIN:
		responseMsg = fmt.Sprintf("Okay that's too short, you gotta be a bit more strict.\nMinimum time is %v.\n", config.minTime)
	case ERR_TIME_ABOVE_MAX:
		responseMsg = fmt.Sprintf("Way too long!\nMaximum time is %v.\n", config.maxTime)
	case ERR_USER_ALREADY_IN_JAIL:
		responseMsg = fmt.Sprintf("User <@!%s> is already in jail! Let's calm down.\n", selectedUser.Username)
	case ERR_COLOR_INVALID_FORMAT:
		responseMsg = fmt.Sprintf("Not a valid color format. Defaulting to `%s`.\n", config.embedColor.hexValue)
	default:
		responseMsg = fmt.Sprintln("Something went wrong, sorry :(")
	}

	return responseMsg
}

func sendErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, responseMsg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: responseMsg,
		},
	})
}

func sendErrorFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, responseMsg string) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: responseMsg,
	})
}
