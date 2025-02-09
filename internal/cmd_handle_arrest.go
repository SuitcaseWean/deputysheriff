package deputysheriff

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CommandsHandlers["arrest"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		selectedUser := i.ApplicationCommandData().Options[0].UserValue(s)
		currArrestUserID = selectedUser.ID

		if config.arrestRole == nil {
			msg := fmt.Sprintln("Please set an arrest role. You can do it by using **/arrest-config-set**") + fmt.Sprintln("If you don't have access to this command, you can kindly ask your Admin people to do it :)")
			sendErrorResponse(s, i, msg)
			return
		}

		if config.annoucementsChannel == nil {
			ch, err := s.Channel(i.ChannelID)
			if err != nil {
				log.Println(err)
				sendErrorResponse(s, i, errMsg(ERR_DEFAULT, nil))
			}
			config.annoucementsChannel = ch
		}

		if errCode, err := createNewArrest(s, i, selectedUser); errCode != 0 {
			log.Println(err)
			sendErrorResponse(s, i, errMsg(errCode, selectedUser))
			return
		}

		reasonValue := ""
		log.Println(arrests[currArrestUserID].reason)
		if (!arrests[currArrestUserID].success) && (arrests[currArrestUserID].reason != "") {
			reasonValue = arrests[currArrestUserID].reason
		}

		err := s.InteractionRespond(i.Interaction, createArrestModal(reasonValue))
		if err != nil {
			log.Println(err)
			sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		}
	}
}

func createArrestModal(reasonValue string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "arrest-modal",
			Title:    "New Arrest",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "reason",
							Label:     "Why is this person being arrested?",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 2000,
							Value:     reasonValue,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "time",
							Label:       "For how long? (default 1min)",
							Style:       discordgo.TextInputShort,
							Placeholder: "[min]m[sec]s (e.g. 1m30s / 3m / 30s)",
							Required:    false,
							MinLength:   2,
							MaxLength:   6,
						},
					},
				},
			},
		},
	}
}

func createNewArrest(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUser *discordgo.User) (int, error) {
	// selectedUser.
	m, err := s.GuildMember(i.GuildID, selectedUser.ID)
	if err != nil {
		return ERR_DEFAULT, err
	}
	if slices.Contains(m.Roles, config.arrestRole.ID) {
		return ERR_USER_ALREADY_IN_JAIL, errors.New("user already arrested")
	}
	if _, ok := arrests[selectedUser.ID]; ok {
		// User already exists but there was a wrong value passed
		return 0, nil
	}
	// UserID doesn't exist -> create a new one
	arrests[selectedUser.ID] = &Arrest{user: selectedUser, reason: "", timeString: config.defaultTime, success: false}
	return 0, nil
}
