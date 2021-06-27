package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/core"
	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

const (
	EmojiCreateReminder    = "⏰"
	EmojiCreateReminderAlt = "⏲️"

	EmojiRefreshDuration  = "🔄"
	EmojiIncreaseDuration = "🔼"
	EmojiDecreaseDuration = "🔽"
	EmojiDeleteReminder   = "🗑️"

	EmojiSuccess = "✅"
	EmojiError   = "❌"
)

const (
	MaximumNoteLength               = 240
	MaximumNumberOfRemindersPerUser = 50

	MinimumReminderDuration = time.Minute
	MaximumReminderDuration = 180 * 24 * time.Hour
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
		log.Printf("[main][HandleMessage] command=%s; query=%s", command, query)
		switch command {
		case "remindme", "remind", "reminder", "help":
			HandleRemindMe(bot, message, query)
		}
	}
}

func HandleRemindMe(bot *discordgo.Session, message *discordgo.MessageCreate, query string) {
	if len(query) == 0 {
		_, _ = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("**Usage:**\n```%sRemindMe <DURATION> [NOTE]```**Where:**\n- `<DURATION>` must be use one of the following formats: `30m`, `6h30m`, `48h`\n- `[NOTE]` is an optional note to attach to the reminder with less than %d characters\n:exclamation: _You can also create a reminder by reacting with %s or %s to a message_", cfg.CommandPrefix, MaximumNoteLength, EmojiCreateReminder, EmojiCreateReminderAlt), message.Reference())
		return
	}
	// Validate duration
	durationArgument := strings.Split(query, " ")[0]
	duration, err := time.ParseDuration(durationArgument)
	if err != nil {
		log.Printf("[main][HandleRemindMe] Failed to parse duration '%s' for %s: %s", query, message.Author.String(), err.Error())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, "Invalid duration format. Try something like `45m`, `1h30m`, or `13h`.", message.Reference())
		if err != nil {
			log.Printf("[main][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	if duration < MinimumReminderDuration || duration > MaximumReminderDuration {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("Duration must between %s and %s", MinimumReminderDuration, MaximumReminderDuration), message.Reference())
		if err != nil {
			log.Printf("[main][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	// Validate note
	note := strings.TrimSpace(strings.TrimPrefix(query, durationArgument))
	if len(note) > MaximumNoteLength {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("Note must have less than %d characters", MaximumNoteLength), message.Reference())
		if err != nil {
			log.Printf("[main][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	// Create the reminder
	_, err = createReminder(bot, message.Author.ID, message.GuildID, message.ChannelID, message.ID, note, time.Now().Add(duration))
	if err != nil {
		log.Printf("[main][HandleRemindMe] Failed to create reminder: %s", err.Error())
		_, err = bot.ChannelMessageSendReply(message.ChannelID, "Error: "+err.Error(), message.Reference())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		return
	}
	_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiSuccess)
}

func HandleReactionAdd(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if reaction.UserID == bot.State.User.ID {
		return
	}
	if reaction.Emoji.Name == EmojiCreateReminder || reaction.Emoji.Name == EmojiCreateReminderAlt {
		_, err := createReminder(bot, reaction.UserID, reaction.GuildID, reaction.ChannelID, reaction.MessageID, "", time.Now().Add(8*time.Hour))
		if err != nil {
			log.Printf("[main][HandleReactionAdd] Failed to create reminder: %s", err.Error())
			return
		}
	}
	// If the user wants to increase the duration
	if reaction.Emoji.Name == EmojiIncreaseDuration || reaction.Emoji.Name == EmojiDecreaseDuration || reaction.Emoji.Name == EmojiDeleteReminder || reaction.Emoji.Name == EmojiRefreshDuration {
		message, err := bot.ChannelMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			return
		}
		// Make sure that the reaction is on one of the bot's messages
		if message.Author.ID == bot.State.User.ID {
			// get the notification message
			reminder, err := database.GetReminderByNotificationMessageID(message.ID)
			if err != nil {
				log.Println("[HandleReactionAdd] Failed to retrieve reminder by notification message id:", err.Error())
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
}

func createReminder(bot *discordgo.Session, userID, guildID, channelID, messageID, note string, when time.Time) (*core.Reminder, error) {
	if len(note) > MaximumNoteLength {
		return nil, fmt.Errorf("note exceeded maximum length of %d", MaximumNoteLength)
	}
	numberOfReminders, _ := database.CountRemindersByUserID(userID)
	if numberOfReminders > MaximumNumberOfRemindersPerUser {
		return nil, fmt.Errorf("you have reached the maximum number of reminders a single user can have (%d)", MaximumNumberOfRemindersPerUser)
	}
	reminder := &core.Reminder{
		UserID:      userID,
		MessageLink: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID),
		Note:        note,
		Time:        when,
	}
	directMessage, err := bot.UserChannelCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create DM with %s: %s", userID, err.Error())
	}
	botMessage, err := bot.ChannelMessageSend(directMessage.ID, reminder.GenerateNotificationMessageContent())
	if err != nil {
		return nil, fmt.Errorf("failed to send DM to %s: %s", userID, err.Error())
	}
	reminder.NotificationMessageID = botMessage.ID
	err = database.CreateReminder(reminder)
	if err != nil {
		return nil, fmt.Errorf("failed to create reminder in database: %s", err.Error())
	}
	_ = bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, EmojiRefreshDuration)
	_ = bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, EmojiIncreaseDuration)
	_ = bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, EmojiDecreaseDuration)
	_ = bot.MessageReactionAdd(botMessage.ChannelID, botMessage.ID, EmojiDeleteReminder)
	return reminder, nil
}
