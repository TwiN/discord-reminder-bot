package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleMessage(bot *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot || message.Author.ID == bot.State.User.ID {
		return
	}
	if strings.HasPrefix(message.Content, botMention) {
		// If the user mentions the bot, we assume they meant to type the command prefix followed by RemindMe
		message.Content = strings.Replace(message.Content, botMention, cfg.CommandPrefix+"remindme", 1)
	}
	if strings.HasPrefix(message.Content, cfg.CommandPrefix) {
		command := strings.Replace(strings.Split(message.Content, " ")[0], cfg.CommandPrefix, "", 1)
		query := strings.TrimSpace(strings.Replace(message.Content, cfg.CommandPrefix+command, "", 1))
		command = strings.ToLower(command)
		log.Printf("command=%s; query=%s", command, query)
		switch command {
		case "remindme":
			HandleRemindMe(bot, message, query)
		case "help":
			HandleHelp(bot, message)
		}
	}
}

func HandleRemindMe(bot *discordgo.Session, message *discordgo.MessageCreate, query string) {
	if len(query) == 0 {
		_, _ = bot.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Usage: ```%sRemindMe <duration>```where `<duration>` must be use one of the following format: `30m`, `6h30m`, `48h`", cfg.CommandPrefix))
		return
	}
	duration, err := time.ParseDuration(query)
	if err != nil {
		log.Printf("[HandleRemindMe] Failed to parse duration '%s' for %s: %s", query, message.Author.String(), err.Error())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, "❌")
		_, err = bot.ChannelMessageSendReply(message.ChannelID, "Invalid duration format. Try something like `45m`, `1h30m`, or `13h`.", message.Reference())
		if err != nil {
			log.Printf("[HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	directMessage, err := bot.UserChannelCreate(message.Author.ID)
	if err != nil {
		log.Printf("[HandleRemindMe] Failed to create DM to %s: %s", message.Author.String(), err.Error())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, "❌")
		return
	} else {
		messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", message.GuildID, message.ChannelID, message.ID)
		_, err = bot.ChannelMessageSend(directMessage.ID, fmt.Sprintf("I will remind you about %s in %s", messageLink, duration.String()))
		if err != nil {
			_ = bot.MessageReactionAdd(message.ChannelID, message.ID, "❌")
			log.Printf("[HandleRemindMe] Failed to send DM to %s: %s", message.Author.String(), err.Error())
			return
		}
	}
	_ = bot.MessageReactionAdd(message.ChannelID, message.ID, "✅")
}

func HandleHelp(bot *discordgo.Session, message *discordgo.MessageCreate) {
	_, _ = bot.ChannelMessageSend(message.ChannelID, fmt.Sprintf(`
__**Commands**__
**%sRemindMe**: Configures a reminder for a specific message

Bugs to report? Create an issue at https://github.com/TwinProduction/discord-reminder-bot`, cfg.CommandPrefix))
}
