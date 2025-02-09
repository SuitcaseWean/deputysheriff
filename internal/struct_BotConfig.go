package deputysheriff

import "github.com/bwmarrin/discordgo"

type BotConfig struct {
	annoucementsChannel *discordgo.Channel
	arrestRole          *discordgo.Role
	defaultTime         string
	minTime             string
	maxTime             string
	embedColor          Color
}