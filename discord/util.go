package discord

import (
	"github.com/bwmarrin/discordgo"
)

func generateMessageEmbed(title, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       title,
		Description: description,
		Color:       color,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: botAvatar},
	}
}
