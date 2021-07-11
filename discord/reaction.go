package discord

import (
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

func HandleReactionAdd(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	handleReaction(bot, reaction.MessageReaction, false)
}

func HandleReactionRemove(bot *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	handleReaction(bot, reaction.MessageReaction, true)
}

func handleReaction(bot *discordgo.Session, reaction *discordgo.MessageReaction, remove bool) {
	// Make sure that the reaction is on one of the bot's messages
	if reaction.UserID == bot.State.User.ID {
		return
	}
	switch reaction.Emoji.Name {
	case EmojiCreateReminder, EmojiCreateReminderAlt:
		// Create a new reminder when a user reacts with EmojiCreateReminder or EmojiCreateReminderAlt
		if !remove {
			handleReactionCreateReminder(bot, reaction)
		}
	case EmojiPageOne, EmojiPageTwo, EmojiPageThree, EmojiPageFour, EmojiPageFive:
		// Navigate page of reminders
		handleReactionListReminders(bot, reaction)
	case EmojiIncreaseDuration, EmojiDecreaseDuration, EmojiDeleteReminder, EmojiRefreshDuration:
		// Modify an existing reminder
		handleReactionModifyReminder(bot, reaction)
	}
}

func handleReactionCreateReminder(bot *discordgo.Session, reaction *discordgo.MessageReaction) {
	_, err := createReminder(bot, reaction.UserID, reaction.GuildID, reaction.ChannelID, reaction.MessageID, "", time.Now().Add(8*time.Hour))
	if err != nil {
		log.Printf("[discord][HandleReactionAdd] Failed to create reminder: %s", err.Error())
		return
	}
}

func handleReactionListReminders(bot *discordgo.Session, reaction *discordgo.MessageReaction) {
	switch reaction.Emoji.Name {
	case EmojiPageOne:
		updateExistingReminderListMessage(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, 1)
	case EmojiPageTwo:
		updateExistingReminderListMessage(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, 2)
	case EmojiPageThree:
		updateExistingReminderListMessage(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, 3)
	case EmojiPageFour:
		updateExistingReminderListMessage(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, 4)
	case EmojiPageFive:
		updateExistingReminderListMessage(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, 5)
	default:
		return // not supported
	}
}

func handleReactionModifyReminder(bot *discordgo.Session, reaction *discordgo.MessageReaction) {
	// get the notification message
	reminder, err := database.GetReminderByNotificationMessageID(reaction.MessageID)
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
		_, _ = updateExistingMessage(bot, reaction.ChannelID, reaction.MessageID, "", reminder.GenerateNotificationMessageContent())
	case EmojiDecreaseDuration:
		// No need to worry about checking if we're under 0, just let it be naturally handled
		reminder.Time = reminder.Time.Add(-time.Hour)
		_ = database.UpdateReminder(reminder)
		_, _ = updateExistingMessage(bot, reaction.ChannelID, reaction.MessageID, "", reminder.GenerateNotificationMessageContent())
	case EmojiRefreshDuration:
		_, _ = updateExistingMessage(bot, reaction.ChannelID, reaction.MessageID, "", reminder.GenerateNotificationMessageContent())
	case EmojiDeleteReminder:
		deleteReminder(bot, reaction.ChannelID, reminder)
	default:
		return // not supported
	}
}

func updateExistingReminderListMessage(bot *discordgo.Session, userID, channelID, messageID string, page int) {
	msg, _, err := createReminderListMessageEmbed(channelID, userID, page)
	if err != nil {
		return
	}
	bot.ChannelMessageEditEmbed(channelID, messageID, msg)
}
