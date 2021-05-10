package database

import (
	"testing"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/core"
)

func TestCreateReminder(t *testing.T) {
	Initialize("sqlite", t.TempDir()+"/test.db")
	db := connect()
	defer db.Close()
	err := CreateReminder(&core.Reminder{
		NotificationMessageID: "1",
		UserID:                "2",
		MessageLink:           "3",
		Note:                  "4",
		Time:                  time.Now().Round(time.Minute),
	})
	if err != nil {
		t.Fatal("failed to create reminder:", err.Error())
	}
	reminder, err := GetReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to retrieve reminder by notification message id:", err.Error())
	}
	if reminder == nil {
		t.Fatal("couldn't find reminder by notification message id from the database")
	}
	if reminder.NotificationMessageID != "1" {
		t.Fatal("NotificationMessageID should've been 1, got", reminder.NotificationMessageID)
	}
	if reminder.UserID != "2" {
		t.Fatal("UserID should've been 2, got", reminder.UserID)
	}
	if reminder.MessageLink != "3" {
		t.Fatal("MessageLink should've been 3, got", reminder.MessageLink)
	}
	if reminder.Note != "4" {
		t.Fatal("Note should've been 4, got", reminder.Note)
	}
	if reminder.Time != time.Now().Round(time.Minute) {
		t.Fatalf("Note should've been %s, got %s", time.Now().Round(time.Minute), reminder.Time)
	}
}

func TestDeleteReminderByNotificationMessageID(t *testing.T) {
	Initialize("sqlite", t.TempDir()+"/test.db")
	db := connect()
	defer db.Close()
	err := CreateReminder(&core.Reminder{
		NotificationMessageID: "1",
		UserID:                "2",
		MessageLink:           "3",
		Note:                  "4",
		Time:                  time.Now().Round(time.Minute),
	})
	if err != nil {
		t.Fatal("failed to create reminder:", err.Error())
	}
	reminder, err := GetReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to retrieve reminder by notification message id:", err.Error())
	}
	if reminder == nil {
		t.Fatal("couldn't find reminder by notification message id from the database")
	}
	err = DeleteReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to delete reminder by notification message id:", err.Error())
	}
	reminder, err = GetReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to retrieve reminder by notification message id:", err.Error())
	}
	if reminder != nil {
		t.Fatal("reminder should've been nil because it was deleted")
	}
}

func TestUpdateReminder(t *testing.T) {
	Initialize("sqlite", t.TempDir()+"/test.db")
	db := connect()
	defer db.Close()
	err := CreateReminder(&core.Reminder{
		NotificationMessageID: "1",
		UserID:                "2",
		MessageLink:           "3",
		Note:                  "4",
		Time:                  time.Now().Round(time.Minute),
	})
	if err != nil {
		t.Fatal("failed to create reminder:", err.Error())
	}
	reminder, err := GetReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to retrieve reminder by notification message id:", err.Error())
	}
	if reminder == nil {
		t.Fatal("couldn't find reminder by notification message id from the database")
	}
	reminder.Note = "updated-note"
	reminder.Time = reminder.Time.Add(time.Hour)
	err = UpdateReminder(reminder)
	if err != nil {
		t.Fatal("failed to delete reminder by notification message id:", err.Error())
	}
	reminder, err = GetReminderByNotificationMessageID("1")
	if err != nil {
		t.Fatal("failed to retrieve reminder by notification message id:", err.Error())
	}
	if reminder == nil {
		t.Fatal("reminder should've existed")
	}
	if reminder.Note != "updated-note" {
		t.Fatal("Note should've been updated-note, got", reminder.Note)
	}
	if reminder.Time != time.Now().Round(time.Minute).Add(time.Hour) {
		t.Fatalf("Note should've been %s, got %s", time.Now().Round(time.Minute).Add(time.Hour), reminder.Time)
	}
}

func TestGetOverdueReminders(t *testing.T) {
	Initialize("sqlite", t.TempDir()+"/test.db")
	db := connect()
	defer db.Close()
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "1", Time: time.Now().Round(time.Minute).Add(time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "2", Time: time.Now().Round(time.Minute).Add(-time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "3", Time: time.Now().Round(time.Minute).Add(3 * time.Hour)})
	overdueReminders, err := GetOverdueReminders()
	if err != nil {
		t.Fatal("failed to retrieve overdue reminders:", err.Error())
	}
	if len(overdueReminders) != 1 {
		t.Fatal("1 reminder should've been overdue, got", len(overdueReminders))
	}
	if overdueReminders[0].NotificationMessageID != "2" {
		t.Fatal("overdue reminder should've been the one with NotificationMessageID 2, got", overdueReminders[0].NotificationMessageID)
	}
}

func TestGetOverdueRemindersRetrievesTheOldestOnesFirst(t *testing.T) {
	Initialize("sqlite", t.TempDir()+"/test.db")
	db := connect()
	defer db.Close()
	now := time.Now().Round(time.Minute)
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "1", Time: now.Add(time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "2", Time: now.Add(-3 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "3", Time: now.Add(-1 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "4", Time: now.Add(-3 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "5", Time: now.Add(-3 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "6", Time: now.Add(-4 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "7", Time: now.Add(-3 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "8", Time: now.Add(2 * time.Hour)})
	_ = CreateReminder(&core.Reminder{NotificationMessageID: "9", Time: now.Add(-1 * time.Hour)})
	overdueReminders, err := GetOverdueReminders()
	if err != nil {
		t.Fatal("failed to retrieve overdue reminders:", err.Error())
	}
	if len(overdueReminders) != 5 {
		t.Fatal("5 reminders should've been overdue, got", len(overdueReminders))
	}
	for _, overdueReminder := range overdueReminders {
		if overdueReminder.Time != now.Add(-3*time.Hour) && overdueReminder.Time != now.Add(-4*time.Hour) {
			t.Fatal("GetOverdueReminders should've only returned the 5 most overdue reminders")
		}
	}
}
