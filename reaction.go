package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var reminder *Reminder // temporary. will be retrieved from database later on.

const (
	IncreaseDuration = "➕"
	DecreaseDuration = "➖"
	DeleteReminder   = "🗑️"
	RefreshDuration  = "🔄"
)

func HandleReactionAdd(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if reaction.UserID == bot.State.User.ID {
		return
	}
	if reaction.Emoji.Name == "⏲️" {
		directMessage, err := bot.UserChannelCreate(reaction.UserID)
		if err != nil {
			log.Printf("[HandleReactionAdd] Failed to create DM to %s: %s", reaction.UserID, err.Error())
			return
		} else {
			messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", reaction.GuildID, reaction.ChannelID, reaction.MessageID)
			botMessage, err := bot.ChannelMessageSend(directMessage.ID, fmt.Sprintf("I will remind you about %s in %s", messageLink, "8h"))
			if err != nil {
				log.Printf("[HandleReactionAdd] Failed to send DM to %s: %s", reaction.UserID, err.Error())
				return
			}
			bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, RefreshDuration)
			bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, IncreaseDuration)
			bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, DecreaseDuration)
			bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, DeleteReminder)
			reminder = &Reminder{
				UserID:                reaction.UserID,
				MessageLink:           messageLink,
				NotificationMessageID: botMessage.ChannelID,
				Time:                  time.Now().Add(8 * time.Hour),
			}
		}
	}
	// If the user wants to increase the duration
	if reaction.Emoji.Name == IncreaseDuration || reaction.Emoji.Name == DecreaseDuration || reaction.Emoji.Name == DeleteReminder || reaction.Emoji.Name == RefreshDuration {
		message, err := bot.ChannelMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			return
		}
		// Make sure that the reaction is on one of the bot's messages
		if message.Author.ID == bot.State.User.ID {
			// get the notification message
			reminder := getReminderByNotificationMessageID(message.ID)
			if reminder == nil {
				// doesn't exist, so we'll ignore the message.
				return
			}
			switch reaction.Emoji.Name {
			case IncreaseDuration:
				reminder.Time = reminder.Time.Add(time.Hour)
				updateReminderInDatabase(reminder)
				bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateMessageContent())
			case DecreaseDuration:
				// XXX: no need to worry about checking if we're under 0, just let it be naturally handled
				reminder.Time = reminder.Time.Add(-time.Hour)
				updateReminderInDatabase(reminder)
				bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateMessageContent())
			case RefreshDuration:
				bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, reminder.GenerateMessageContent())
			case DeleteReminder:
				bot.ChannelMessageEdit(reaction.ChannelID, reaction.MessageID, "~~"+reminder.GenerateMessageContent()+"~~")
			default:
				return // not supported
			}
		}
	}
}

func getReminderByNotificationMessageID(messageID string) *Reminder {
	// TODO: get reminder by NotificationMessageID from the database, but temporarily going to directly use the reminder object
	return reminder
}

func updateReminderInDatabase(reminder *Reminder) {

}
