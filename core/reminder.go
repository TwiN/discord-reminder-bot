package core

import (
	"time"
)

type Reminder struct {
	NotificationMessageID string    // ID of the message used to manage the reminder
	UserID                string    // ID of the user to notify
	MessageLink           string    // Link to the message that the user wants to be reminded about
	Note                  string    // Note attached to the reminder
	Time                  time.Time // Time at which the reminder is due for
}

func (r Reminder) GenerateNotificationMessageContent() string {
	if time.Until(r.Time) < 0 {
		return "I will remind you about " + r.MessageLink + " at " + r.Time.Format(time.RFC3339)
	}
	return "I will remind you about " + r.MessageLink + " in " + time.Until(r.Time).Round(time.Second).String()
}

func (r Reminder) GenerateReminderMessageContent() string {
	if len(r.Note) > 0 {
		return "You asked me to remind you about " + r.MessageLink + " and attached the following note:\n```" + r.Note + "```"
	}
	return "You asked me to remind you about this message: " + r.MessageLink
}
