package deputysheriff

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Embed struct{}

func (e *Embed) arrest(i *discordgo.InteractionCreate, a Arrest, t int64) discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "What happened?",
			Value: fmt.Sprintf("**<@!%s> was arrested by <@!%s>! 😱**", a.user.ID, i.Interaction.Member.User.ID),
		},
		{
			Name:  "Reason",
			Value: a.reason,
		},
	}
	// if t != "0s" {
	// }
	timeField := discordgo.MessageEmbedField{
		Name:  "You will see them again in ",
		Value: fmt.Sprintf("<t:%v:R>", t),
	}
	fields = append(fields, &timeField)
	return discordgo.MessageEmbed{
		Title:  "Sheriff report",
		Fields: fields,
		Color:  config.embedColor.intValue,
	}
}
func (e *Embed) release(a Arrest) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       "Sheriff report",
		Description: fmt.Sprintf("**<@!%s> was released from jail!**", a.user.ID),
		Color:       config.embedColor.intValue,
	}
}
