package discord

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
		message.Content = strings.Replace(message.Content, botMention, botCommandPrefix+"remindme", 1)
	}
	if strings.HasPrefix(message.Content, botCommandPrefix) {
		command := strings.Replace(strings.Split(message.Content, " ")[0], botCommandPrefix, "", 1)
		query := strings.TrimSpace(strings.Replace(message.Content, botCommandPrefix+command, "", 1))
		command = strings.ToLower(command)
		log.Printf("[discord][HandleMessage] command=%s; arguments=%s", command, query)
		switch command {
		case "remindme", "remind", "reminder", "help":
			HandleRemindMe(bot, message, query)
		case "list", "reminders", "view":
			HandleListReminders(bot, message)
		}
	}
}

func HandleRemindMe(bot *discordgo.Session, message *discordgo.MessageCreate, query string) {
	if len(query) == 0 {
		_, _ = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("**Usage:**\n```%sRemindMe <DURATION> [NOTE]```**Where:**\n- `<DURATION>` must be use one of the following formats: `30m`, `6h30m`, `48h`\n- `[NOTE]` is an optional note to attach to the reminder with less than %d characters\n\n:information_source: _You can also create a reminder by reacting with %s, %s or %s to a message, and you can view your reminders by using `%slist`._", botCommandPrefix, MaximumNoteLength, EmojiCreateReminder, EmojiCreateReminderAlt1, EmojiCreateReminderAlt2, botCommandPrefix), message.Reference())
		return
	}
	// Validate duration
	durationArgument := strings.Split(query, " ")[0]
	duration, err := time.ParseDuration(durationArgument)
	if err != nil {
		log.Printf("[discord][HandleRemindMe] Failed to parse duration '%s' for %s: %s", query, message.Author.String(), err.Error())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("Invalid duration format: ```%s```\n:information_source: _Try something like `45m`, `1h30m`, or `13h`._", err.Error()), message.Reference())
		if err != nil {
			log.Printf("[discord][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	if duration < MinimumReminderDuration || duration > MaximumReminderDuration {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("Duration must between %s and %s", MinimumReminderDuration, MaximumReminderDuration), message.Reference())
		if err != nil {
			log.Printf("[discord][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	// Validate note
	note := strings.TrimSpace(strings.TrimPrefix(query, durationArgument))
	if len(note) > MaximumNoteLength {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		_, err = bot.ChannelMessageSendReply(message.ChannelID, fmt.Sprintf("Note must have less than %d characters", MaximumNoteLength), message.Reference())
		if err != nil {
			log.Printf("[discord][HandleRemindMe] Failed to reply to message: %s", err.Error())
		}
		return
	}
	// Create the reminder
	_, err = createReminder(bot, message.Author.ID, message.GuildID, message.ChannelID, message.ID, note, time.Now().Add(duration))
	if err != nil {
		log.Printf("[discord][HandleRemindMe] Failed to create reminder: %s", err.Error())
		_, err = bot.ChannelMessageSendReply(message.ChannelID, "Error: "+err.Error(), message.Reference())
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		return
	}
	_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiSuccess)
}

func HandleListReminders(bot *discordgo.Session, message *discordgo.MessageCreate) {
	directMessageChannel, err := bot.UserChannelCreate(message.Author.ID)
	if err != nil {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		log.Println("[discord][handleReminderListPage] Failed to open direct message:", err.Error())
		return
	}
	msg, numberOfPages, err := createReminderListMessageEmbed(directMessageChannel.ID, message.Author.ID, 1)
	if err != nil {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		log.Println("[discord][HandleListReminders] Failed to create embed message:", err.Error())
		return
	}
	messageSent, err := bot.ChannelMessageSendEmbed(directMessageChannel.ID, msg)
	if err != nil {
		_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiError)
		log.Println("[discord][HandleListReminders] Failed to send message:", err.Error())
		return
	}
	_ = bot.MessageReactionAdd(message.ChannelID, message.ID, EmojiSuccess)
	_ = bot.MessageReactionAdd(messageSent.ChannelID, messageSent.ID, EmojiPageOne)
	if numberOfPages >= 2 {
		_ = bot.MessageReactionAdd(messageSent.ChannelID, messageSent.ID, EmojiPageTwo)
	}
	if numberOfPages >= 3 {
		_ = bot.MessageReactionAdd(messageSent.ChannelID, messageSent.ID, EmojiPageThree)
	}
	if numberOfPages >= 4 {
		_ = bot.MessageReactionAdd(messageSent.ChannelID, messageSent.ID, EmojiPageFour)
	}
	if numberOfPages >= 5 {
		_ = bot.MessageReactionAdd(messageSent.ChannelID, messageSent.ID, EmojiPageFive)
	}
}
