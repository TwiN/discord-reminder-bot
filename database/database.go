package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/core"
	_ "modernc.org/sqlite"
)

var (
	databaseDriver string
	databasePath   string
)

func Initialize(driver, path string) {
	databaseDriver = driver
	databasePath = path
	log.Printf("[database][Initialize] Beginning schema migration on database with driver=%s", driver)
	createSchema()
}

func connect() *sql.DB {
	db, err := sql.Open(databaseDriver, databasePath)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// createSchema creates the schema required to perform all database operations.
func createSchema() {
	db := connect()
	defer db.Close()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS reminder (
			notification_message_id VARCHAR(64) PRIMARY KEY, 
			user_id                 VARCHAR(64), 
			message_link            VARCHAR(128),
			note                    VARCHAR(255),
			reminder_time           TIMESTAMP
		)
	`)
	if err != nil {
		panic("unable to create table in database: " + err.Error())
	}
}

func CreateReminder(reminder *core.Reminder) error {
	start := time.Now()
	db := connect()
	defer db.Close()
	_, err := db.Exec(
		"INSERT INTO reminder (notification_message_id, user_id, message_link, note, reminder_time) VALUES ($1, $2, $3, $4, $5)",
		reminder.NotificationMessageID,
		reminder.UserID,
		reminder.MessageLink,
		reminder.Note,
		reminder.Time,
	)
	if err != nil {
		log.Printf("[database][CreateReminder] Failed to create reminder for NotificationMessageID=%s; duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][CreateReminder] Created reminder for NotificationMessageID=%s in duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	}
	return err
}

func GetReminderByNotificationMessageID(messageID string) (*core.Reminder, error) {
	start := time.Now()
	db := connect()
	defer db.Close()
	rows, err := db.Query("SELECT notification_message_id, user_id, message_link, note, reminder_time FROM reminder WHERE notification_message_id = $1", messageID)
	if err != nil {
		return nil, err
	}
	var reminder *core.Reminder
	for rows.Next() {
		reminder = &core.Reminder{}
		_ = rows.Scan(&reminder.NotificationMessageID, &reminder.UserID, &reminder.MessageLink, &reminder.Note, &reminder.Time)
		break
	}
	_ = rows.Close()
	if reminder == nil {
		log.Printf("[database][GetReminderByNotificationMessageID] No reminder for NotificationMessageID=%s found; duration=%dms", messageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][GetReminderByNotificationMessageID] Got reminder for NotificationMessageID=%s in duration=%dms", messageID, time.Since(start).Milliseconds())
	}
	return reminder, nil
}

func UpdateReminderInDatabase(reminder *core.Reminder) error {
	start := time.Now()
	db := connect()
	defer db.Close()
	_, err := db.Exec("UPDATE reminder SET reminder_time = $1, note = $2 WHERE notification_message_id = $3", reminder.Time, reminder.Note, reminder.NotificationMessageID)
	if err != nil {
		log.Printf("[database][UpdateReminderInDatabase] Failed to update reminder with NotificationMessageID=%s; duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][UpdateReminderInDatabase] Updated reminder with NotificationMessageID=%s in duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	}
	return err
}

func GetExpiredReminders() ([]*core.Reminder, error) {
	start := time.Now()
	db := connect()
	defer db.Close()
	rows, err := db.Query("SELECT notification_message_id, user_id, message_link, note, reminder_time FROM reminder WHERE reminder_time < $1 LIMIT 5", time.Now())
	if err != nil {
		return nil, err
	}
	var reminders []*core.Reminder
	for rows.Next() {
		reminder := &core.Reminder{}
		_ = rows.Scan(&reminder.NotificationMessageID, &reminder.UserID, &reminder.MessageLink, &reminder.Note, &reminder.Time)
		reminders = append(reminders, reminder)
	}
	_ = rows.Close()
	if len(reminders) > 0 {
		log.Printf("[database][GetExpiredReminders] Got %d reminders in duration=%dms", len(reminders), time.Since(start).Milliseconds())
	}
	return reminders, nil
}

func DeleteReminderByNotificationMessageID(messageID string) error {
	start := time.Now()
	db := connect()
	defer db.Close()
	_, err := db.Exec("DELETE FROM reminder WHERE notification_message_id = $1", messageID)
	if err != nil {
		log.Printf("[database][DeleteReminderByNotificationMessageID] Failed to delete reminder with NotificationMessageID=%s; duration=%dms", messageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][DeleteReminderByNotificationMessageID] Deleted reminder with NotificationMessageID=%s in duration=%dms", messageID, time.Since(start).Milliseconds())
	}
	return err
}
