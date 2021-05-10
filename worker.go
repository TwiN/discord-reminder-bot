package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/config"
	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

func worker(bot *discordgo.Session) {
	for {
		time.Sleep(10 * time.Second)
		reminders, err := database.GetOverdueReminders()
		if err != nil {
			// TODO: if errors 5 times in a row, panic
			log.Println("[main][worker] Failed to retrieve expired reminders from database:", err.Error())
			continue
		}
		if len(reminders) > 0 {
			numberOfQueuedReminders, err := database.CountReminders()
			if err == nil {
				_ = bot.UpdateListeningStatus(fmt.Sprintf("%d queued reminders", numberOfQueuedReminders))
			} else {
				_ = bot.UpdateListeningStatus(config.Get().CommandPrefix + "RemindMe")
			}
		} else {
			_ = bot.UpdateListeningStatus(config.Get().CommandPrefix + "RemindMe")
		}
		for _, reminder := range reminders {
			directMessage, err := bot.UserChannelCreate(reminder.UserID)
			if err != nil {
				log.Printf("[main][worker] Failed to create DM with %s: %s", reminder.UserID, err.Error())
				_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
				continue
			}
			_, err = bot.ChannelMessageSend(directMessage.ID, reminder.GenerateReminderMessageContent())
			if err != nil {
				log.Printf("[main][worker] Failed to send DM to %s: %s", reminder.UserID, err.Error())
				_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
				continue
			}
			// Cross notification message
			_, _ = bot.ChannelMessageEdit(directMessage.ID, reminder.NotificationMessageID, "~~"+reminder.GenerateNotificationMessageContent()+"~~")
			_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
		}
	}
}
