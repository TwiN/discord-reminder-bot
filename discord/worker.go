package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/database"
	"github.com/bwmarrin/discordgo"
)

func worker(bot *discordgo.Session) {
	for {
		time.Sleep(10 * time.Second)
		reminders, err := database.GetOverdueReminders()
		if err != nil {
			// TODO: if errors 5 times in a row, panic
			log.Println("[discord][worker] Failed to retrieve expired reminders from database:", err.Error())
			continue
		}
		if len(reminders) > 0 {
			numberOfQueuedReminders, err := database.CountReminders()
			if err == nil {
				_ = bot.UpdateListeningStatus(fmt.Sprintf("%d queued reminders", numberOfQueuedReminders))
			} else {
				_ = bot.UpdateListeningStatus(botCommandPrefix + "RemindMe")
			}
		} else {
			_ = bot.UpdateListeningStatus(botCommandPrefix + "RemindMe")
		}
		for _, reminder := range reminders {
			directMessage, err := sendDirectMessage(bot, reminder.UserID, "", reminder.GenerateReminderMessageContent())
			if err != nil {
				log.Printf("[discord][worker] Error: %s", err.Error())
				_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
				continue
			}
			// Cross notification message
			_, _ = updateExistingMessage(bot, directMessage.ChannelID, reminder.NotificationMessageID, "", "~~"+reminder.GenerateNotificationMessageContent()+"~~")
			_ = database.DeleteReminderByNotificationMessageID(reminder.NotificationMessageID)
		}
	}
}
