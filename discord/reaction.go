package discord

import (
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

func HandleReactionAdd(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if reaction.UserID == bot.State.User.ID {
		return
	}
	// Create a new reminder
	if reaction.Emoji.Name == EmojiCreateReminder || reaction.Emoji.Name == EmojiCreateReminderAlt {
		handleReactionCreateReminder(bot, reaction)
	}
	// Modify an existing reminder
	if reaction.Emoji.Name == EmojiIncreaseDuration || reaction.Emoji.Name == EmojiDecreaseDuration || reaction.Emoji.Name == EmojiDeleteReminder || reaction.Emoji.Name == EmojiRefreshDuration {
		handleReactionModifyReminder(bot, reaction)
	}
}

func handleReactionCreateReminder(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	_, err := createReminder(bot, reaction.UserID, reaction.GuildID, reaction.ChannelID, reaction.MessageID, "", time.Now().Add(8*time.Hour))
	if err != nil {
		log.Printf("[discord][HandleReactionAdd] Failed to create reminder: %s", err.Error())
		return
	}
}

func handleReactionModifyReminder(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	message, err := bot.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return
	}
	// Make sure that the reaction is on one of the bot's messages
	if message.Author.ID == bot.State.User.ID {
		// get the notification message
		reminder, err := database.GetReminderByNotificationMessageID(message.ID)
		if err != nil {
			log.Println("[discord][HandleReactionAdd] Failed to retrieve reminder by notification message id:", err.Error())
			return
		}
		if reminder == nil {
			// doesn't exist, so we'll ignore the message.
			return
		}
		switch reaction.Emoji.Name {
		case EmojiIncreaseDuration:
			reminder.Time = reminder.Time.Add(time.Hour)
			_ = database.UpdateReminder(reminder)
			_, _ = bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateNotificationMessageContent())
		case EmojiDecreaseDuration:
			// XXX: no need to worry about checking if we're under 0, just let it be naturally handled
			reminder.Time = reminder.Time.Add(-time.Hour)
			_ = database.UpdateReminder(reminder)
			_, _ = bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateNotificationMessageContent())
		case EmojiRefreshDuration:
			_, _ = bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateNotificationMessageContent())
		case EmojiDeleteReminder:
			_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
			_, _ = bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, "~~"+reminder.GenerateNotificationMessageContent()+"~~")
		default:
			return // not supported
		}
	}
}
