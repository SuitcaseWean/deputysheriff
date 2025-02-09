package deputysheriff

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Arrest struct {
	user               *discordgo.User
	success            bool
	timeString, reason string
}

func (a Arrest) ValidateTime() (int, error) {
	t, err := time.ParseDuration(a.timeString)
	if err != nil {
		return ERR_TIME_INVALID_FORMAT, err
	}

	min, _ := time.ParseDuration(config.minTime)
	max, _ := time.ParseDuration(config.maxTime)

	if t.Seconds() < min.Seconds() {
		return ERR_TIME_BELOW_MIN, errors.New("time below minimum")
	}
	if t.Seconds() > max.Seconds() {
		return ERR_TIME_ABOVE_MAX, errors.New("time above maximum")
	}

	return 0, nil
}

func (a *Arrest) makeArrest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.GuildMemberRoleAdd(i.GuildID, a.user.ID, config.arrestRole.ID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintln("Succesful arrest"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	timeStampNow := time.Now().Unix()
	embed := sendEmbed.arrest(i, *a, timeStampNow)
	_, err := s.ChannelMessageSendEmbed(config.annoucementsChannel.ID, &embed)
	if err != nil {
		log.Println(err)
		sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		return
	}

	a.success = true
	log.Println("Role Added succesfully")

	t, _ := time.ParseDuration(a.timeString)
	seconds := int(t.Seconds())
	dur := t

	for range time.Tick(time.Second) {
		// Update time in channel?
		if seconds == 0 {
			break
		}
		seconds--
		dur -= time.Second

		// embed := sendEmbed.arrest(i, *a, dur.String())
		// _, err = s.ChannelMessageEditEmbed(config.annoucementsChannel.ID, msg.ID, &embed)
		// if err != nil {
		// 	log.Println(err)
		// 	sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		// 	return
		// }
	}
	defer a.breakFree(s, i)
}

func (a *Arrest) breakFree(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.GuildMemberRoleRemove(i.GuildID, a.user.ID, config.arrestRole.ID)
	if err != nil {
		log.Println(err)
		sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		return
	}
	delete(arrests, a.user.ID)

	embed := sendEmbed.release(*a)
	_, err = s.ChannelMessageSendEmbed(config.annoucementsChannel.ID, &embed)
	if err != nil {
		log.Println(err)
		sendErrorFollowup(s, i, errMsg(ERR_DEFAULT, nil))
		return
	}

	log.Println("Role Removed succesfully")
}
