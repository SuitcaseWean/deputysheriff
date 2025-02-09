package deputysheriff

import "github.com/bwmarrin/discordgo"

var currArrestUserID string

var (
	CommandsDefinitions = []*discordgo.ApplicationCommand{}
	CommandsHandlers    = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	ComponentsHandlers  = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	arrests = make(map[string]*Arrest) // key == ArrestedUserID
	// defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	sendEmbed = Embed{}

	config = BotConfig{
		annoucementsChannel: nil,
		arrestRole:          nil,
		defaultTime:         "1m",
		minTime:             "30s",
		maxTime:             "5m",
		embedColor:          Color{hexValue: "#000000", intValue: 0},
	}
)
