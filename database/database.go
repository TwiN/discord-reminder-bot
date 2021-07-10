package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/TwinProduction/discord-reminder-bot/core"
	_ "modernc.org/sqlite"
)

var db *sql.DB

// Initialize the database and creates the schema if it doesn't already exist in the file specified
func Initialize(driver, path string) (err error) {
	db, err = sql.Open(driver, path)
	if err != nil {
		return err
	}
	log.Printf("[database][Initialize] Beginning schema migration on database with driver=%s", driver)
	return createSchema()
}

// createSchema creates the schema required to perform all database operations.
func createSchema() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS reminder (
			notification_message_id VARCHAR(64) PRIMARY KEY, 
			user_id                 VARCHAR(64), 
			message_link            VARCHAR(128),
			note                    VARCHAR(255),
			reminder_time           TIMESTAMP,
		    -- If I implement repeating intervals, I need to support keywords like "everyday in (time in 8h)" OR I could allow users to configure their timezones 
		    -- (and persist it in a separate table), AND I need to create a command to print all reminders
		    --
		    -- Something like r!remindme every [DURATION]  [INTERVAL DURATION] [NOTE]
		    -- e.g. "r!remindme 10h every 24h Go jog" would remind somebody "Go jog" in 10 hours from now, every 24 hours
		    --
		    repeating               INTEGER DEFAULT FALSE 
		    
		)
	`)
	return err
}

func CreateReminder(reminder *core.Reminder) error {
	start := time.Now()
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

// GetReminderByNotificationMessageID retrieves a reminder by its NotificationMessageID.
// NotificationMessageID is always unique, because it represents the message ID of the
// message sent to the user by direct message
func GetReminderByNotificationMessageID(messageID string) (*core.Reminder, error) {
	start := time.Now()
	rows, err := db.Query("SELECT rowid, notification_message_id, user_id, message_link, note, reminder_time FROM reminder WHERE notification_message_id = $1", messageID)
	if err != nil {
		return nil, err
	}
	var reminder *core.Reminder
	for rows.Next() {
		reminder = &core.Reminder{}
		_ = rows.Scan(&reminder.ID, &reminder.NotificationMessageID, &reminder.UserID, &reminder.MessageLink, &reminder.Note, &reminder.Time)
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

// UpdateReminder updates a reminder
// Note that the only fields supported for updates are Reminder.Note and Reminder.Time
func UpdateReminder(reminder *core.Reminder) error {
	start := time.Now()
	_, err := db.Exec("UPDATE reminder SET reminder_time = $1, note = $2 WHERE notification_message_id = $3", reminder.Time, reminder.Note, reminder.NotificationMessageID)
	if err != nil {
		log.Printf("[database][UpdateReminder] Failed to update reminder with NotificationMessageID=%s; duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][UpdateReminder] Updated reminder with NotificationMessageID=%s in duration=%dms", reminder.NotificationMessageID, time.Since(start).Milliseconds())
	}
	return err
}

// GetOverdueReminders retrieves at most 5 reminders who have exceeded the time at which said reminder was due
func GetOverdueReminders() ([]*core.Reminder, error) {
	start := time.Now()
	rows, err := db.Query(
		"SELECT rowid, notification_message_id, user_id, message_link, note, reminder_time FROM reminder WHERE reminder_time < $1 ORDER BY reminder_time LIMIT 5",
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	var reminders []*core.Reminder
	for rows.Next() {
		reminder := &core.Reminder{}
		_ = rows.Scan(&reminder.ID, &reminder.NotificationMessageID, &reminder.UserID, &reminder.MessageLink, &reminder.Note, &reminder.Time)
		reminders = append(reminders, reminder)
	}
	_ = rows.Close()
	if len(reminders) > 0 {
		log.Printf("[database][GetOverdueReminders] Got %d reminders in duration=%dms", len(reminders), time.Since(start).Milliseconds())
	}
	return reminders, nil
}

func CountReminders() (int, error) {
	rows, err := db.Query("SELECT COUNT(1) FROM reminder")
	if err != nil {
		return 0, err
	}
	var numberOfReminders int
	for rows.Next() {
		_ = rows.Scan(&numberOfReminders)
		break
	}
	_ = rows.Close()
	return numberOfReminders, nil
}

func GetRemindersByUserID(userId string, page, pageSize int) ([]*core.Reminder, error) {
	rows, err := db.Query("SELECT rowid, notification_message_id, user_id, message_link, note, reminder_time FROM reminder WHERE user_id = $1 ORDER BY reminder_time LIMIT $2 OFFSET $3", userId, pageSize, page*pageSize)
	if err != nil {
		return nil, err
	}
	var reminders []*core.Reminder
	for rows.Next() {
		reminder := &core.Reminder{}
		_ = rows.Scan(&reminder.ID, &reminder.NotificationMessageID, &reminder.UserID, &reminder.MessageLink, &reminder.Note, &reminder.Time)
		reminders = append(reminders, reminder)
	}
	_ = rows.Close()
	return reminders, nil
}

func CountRemindersByUserID(userId string) (int, error) {
	rows, err := db.Query("SELECT COUNT(1) FROM reminder WHERE user_id = $1", userId)
	if err != nil {
		return 0, err
	}
	var numberOfReminders int
	for rows.Next() {
		_ = rows.Scan(&numberOfReminders)
		break
	}
	_ = rows.Close()
	return numberOfReminders, nil
}

func DeleteReminderByNotificationMessageID(messageID string) error {
	start := time.Now()
	_, err := db.Exec("DELETE FROM reminder WHERE notification_message_id = $1", messageID)
	if err != nil {
		log.Printf("[database][DeleteReminderByNotificationMessageID] Failed to delete reminder with NotificationMessageID=%s; duration=%dms", messageID, time.Since(start).Milliseconds())
	} else {
		log.Printf("[database][DeleteReminderByNotificationMessageID] Deleted reminder with NotificationMessageID=%s in duration=%dms", messageID, time.Since(start).Milliseconds())
	}
	return err
}
