package core

import (
	"time"
)

type Reminder struct {
	NotificationMessageID string // MessageID of the message used to manage the reminder
	UserID                string
	MessageLink           string
	Note                  string
	Time                  time.Time
}

func (r Reminder) GenerateNotificationMessageContent() string {
	return "I will remind you about " + r.MessageLink + " in " + time.Until(r.Time).Round(time.Second).String()
}

func (r Reminder) GenerateReminderMessageContent() string {
	if len(r.Note) > 0 {
		return "You asked me to remind you about " + r.MessageLink + " and attached the following note:\n```" + r.Note + "```"
	}
	return "You asked me to remind you about this message: " + r.MessageLink
}
