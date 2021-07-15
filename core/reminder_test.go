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
	if expected := "I will remind you about [this message](<MessageLink>) in 1 hour"; reminder.GenerateNotificationMessageContent() != expected {
		t.Errorf("expected '%s', got '%s'", expected, reminder.GenerateNotificationMessageContent())
	}
	reminder.Time = time.Now().Add(26*time.Hour + 30*time.Minute)
	if expected := "I will remind you about [this message](<MessageLink>) in 1 day, 2 hours and 30 minutes"; reminder.GenerateNotificationMessageContent() != expected {
		t.Errorf("expected '%s', got '%s'", expected, reminder.GenerateNotificationMessageContent())
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
