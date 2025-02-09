package deputysheriff

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	ComponentsHandlers["arrest-modal"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ModalSubmitData()

		a, ok := arrests[currArrestUserID]
		if !ok {
			log.Println("Something went wrong, user not added")
			sendErrorResponse(s, i, errMsg(ERR_DEFAULT, nil))
			return
		}

		a.reason = data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		timeSubmitted := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		if timeSubmitted != "" {
			a.timeString = timeSubmitted
		}

		if errCode, err := a.ValidateTime(); errCode != 0 {
			log.Printf("%v", err)
			a.success = false
			sendErrorResponse(s, i, errMsg(errCode, a.user))
			return
		}

		a.makeArrest(s, i)
	}
}
