package main

import "time"

type Reminder struct {
	UserID                string
	MessageLink           string
	NotificationMessageID string // MessageID of the message used to manage the reminder
	Time                  time.Time
}

func (r Reminder) GenerateMessageContent() string {
	return "I will remind you about " + r.MessageLink + " in " + time.Until(r.Time).String()
}
