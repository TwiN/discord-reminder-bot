package discord

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/config"
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

	EmojiPageOne   = "1️⃣"
	EmojiPageTwo   = "2️⃣"
	EmojiPageThree = "3️⃣"
	EmojiPageFour  = "4️⃣"
	EmojiPageFive  = "5️⃣"

	EmojiSuccess = "✅"
	EmojiError   = "❌"
)

const (
	MaximumNoteLength               = 240
	MaximumNumberOfRemindersPerUser = 50

	MinimumReminderDuration = time.Minute
	MaximumReminderDuration = 180 * 24 * time.Hour

	ReminderListPageSize = 7
)

var (
	botMention       string
	botAvatar        string
	botCommandPrefix string
)

func Start(bot *discordgo.Session, cfg *config.Config) {
	botMention = "<@!" + bot.State.User.ID + ">"
	botAvatar = bot.State.User.AvatarURL("128")
	botCommandPrefix = cfg.CommandPrefix
	bot.AddHandler(HandleMessage)
	bot.AddHandler(HandleReactionAdd)
	bot.AddHandler(HandleReactionRemove)
	_ = bot.UpdateListeningStatus(botCommandPrefix + "RemindMe")
	go worker(bot)
}

// sendDirectMessage sends a direct message to a user and returns the ID of the message sent
func sendDirectMessage(bot *discordgo.Session, userID string, message string) (*discordgo.Message, error) {
	directMessageChannel, err := bot.UserChannelCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create DM with %s: %s", userID, err.Error())
	}
	directMessage, err := bot.ChannelMessageSend(directMessageChannel.ID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send DM to %s: %s", userID, err.Error())
	}
	return directMessage, nil
}

func createReminder(bot *discordgo.Session, userID, guildID, channelID, messageID, note string, when time.Time) (*core.Reminder, error) {
	if len(note) > MaximumNoteLength {
		return nil, fmt.Errorf("note exceeded maximum length of %d", MaximumNoteLength)
	}
	numberOfReminders, _ := database.CountRemindersByUserID(userID)
	if numberOfReminders > MaximumNumberOfRemindersPerUser {
		return nil, fmt.Errorf("you have reached the maximum number of reminders a single user can have (%d)", MaximumNumberOfRemindersPerUser)
	}
	if len(guildID) == 0 {
		guildID = "@me"
	}
	reminder := &core.Reminder{
		UserID:      userID,
		MessageLink: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID),
		Note:        note,
		Time:        when,
	}
	directMessage, err := sendDirectMessage(bot, userID, reminder.GenerateNotificationMessageContent())
	if err != nil {
		return nil, err
	}
	reminder.NotificationMessageID = directMessage.ID
	err = database.CreateReminder(reminder)
	if err != nil {
		return nil, fmt.Errorf("failed to create reminder in database: %s", err.Error())
	}
	_ = bot.MessageReactionAdd(directMessage.ChannelID, directMessage.ID, EmojiRefreshDuration)
	_ = bot.MessageReactionAdd(directMessage.ChannelID, directMessage.ID, EmojiIncreaseDuration)
	_ = bot.MessageReactionAdd(directMessage.ChannelID, directMessage.ID, EmojiDecreaseDuration)
	_ = bot.MessageReactionAdd(directMessage.ChannelID, directMessage.ID, EmojiDeleteReminder)
	return reminder, nil
}

// createReminderListMessageEmbed creates a MessageEmbed and returns the total number of pages
// for the number of reminders that the user has
func createReminderListMessageEmbed(notificationMessageChannelID, userID string, page int) (*discordgo.MessageEmbed, int, error) {
	reminders, err := database.GetRemindersByUserID(userID, page-1, ReminderListPageSize)
	if err != nil {
		return nil, 0, err
	}
	var fields []*discordgo.MessageEmbedField
	for _, reminder := range reminders {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d: In %s from now", reminder.ID, time.Until(reminder.Time).Round(time.Second).String()),
			Value: reminder.GenerateReminderMessageContentInList(notificationMessageChannelID),
		})
	}
	var description string
	if len(fields) == 0 {
		description = "_No reminders to display_"
	} else {
		description = "Below is a list of all reminders on page " + strconv.Itoa(page)
	}
	numberOfReminders, _ := database.CountRemindersByUserID(userID)
	numberOfPages := (numberOfReminders / ReminderListPageSize) + 1
	msg := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Reminders",
		Description: description,
		Fields:      fields,
		Color:       0x20B020,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    botAvatar,
			Width:  128,
			Height: 128,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Page %d out of %d", page, numberOfPages),
			IconURL: botAvatar,
		},
	}
	return msg, numberOfPages, nil
}
