package core

import (
	"testing"
	"time"
)

func TestReminder_GenerateNotificationMessageContent(t *testing.T) {
	reminder := &Reminder{
		NotificationMessageID: "<NotificationMessageID>",
		UserID:                "<UserID>",
		MessageLink:           "<MessageLink>",
		Note:                  "<Note>",
		Time:                  time.Now().Add(time.Hour),
	}
	if reminder.GenerateNotificationMessageContent() != "I will remind you about [this message](<MessageLink>) in 1h0m0s" {
		t.Error("expected 'I will remind you about [this message](<MessageLink>) in 1h0m0s', got", reminder.GenerateNotificationMessageContent())
	}
}

func TestReminder_GenerateReminderMessageContent(t *testing.T) {
	reminder := &Reminder{
		NotificationMessageID: "<NotificationMessageID>",
		UserID:                "<UserID>",
		MessageLink:           "<MessageLink>",
		Note:                  "<Note>",
		Time:                  time.Now().Add(time.Hour),
	}
	if reminder.GenerateReminderMessageContent() != "You asked me to remind you about [this message](<MessageLink>) and attached the following note:\n```<Note>```" {
		t.Error("expected 'You asked me to remind you about [this message](<MessageLink>) and attached the following note:\n```<Note>```', got", reminder.GenerateReminderMessageContent())
	}
}
