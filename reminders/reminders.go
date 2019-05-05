package reminders

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

type Reminderer interface {
	Create(rem Reminder) error
	GetByID(id int) (*Reminder, error)
	DeleteByID(id int) error
	SetSnooze(id int) error
	GetTodayByUser(userID string) ([]Reminder, error)
	DeleteAllByUser(userID string) error
	GetExpired(snoozePeriod string) ([]Reminder, error)
}

type Reminder struct {
	Id               int       `json:"id"`
	Text             string    `json:"text"`
	UserID           string    `json:"user_id"`
	RemindAt         time.Time `json:"remind_at"`
	RemindAtOriginal string    `json:"remind_at_original"`
	Snoozed          bool      `json:"snoozed"`
}

func (r Reminder) String() string {
	b, _ := json.Marshal(r)

	return string(b)
}

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) Create(rem Reminder) error {
	_, err := m.db.Exec(insertReminder, rem.Text, rem.UserID, rem.RemindAt, rem.RemindAtOriginal)

	return err
}

func (m *Manager) DeleteByID(id int) error {
	_, err := m.db.Exec(deleteReminderByID, id)

	return err
}

func (m *Manager) SetSnooze(id int) error {
	_, err := m.db.Exec(setSnoozeReminderByID, id)

	return err
}

func (m *Manager) DeleteAllByUser(userID string) error {
	_, err := m.db.Exec(deleteReminderByUserID, userID)

	return err
}

func (m *Manager) GetByID(id int) (*Reminder, error) {
	var err error
	var rem Reminder

	row, err := m.db.Query(getReminderByID, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(
			&rem.Id,
			&rem.Text,
			&rem.UserID,
			&rem.RemindAt,
			&rem.RemindAtOriginal,
			&rem.Snoozed,
		)
		if err != nil {
			return nil, err
		}

	}

	return &rem, nil
}

func (m *Manager) GetTodayByUser(userID string) ([]Reminder, error) {
	var err error
	var reminders []Reminder

	row, err := m.db.Query(getRemindersTodayByUID, userID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var rem Reminder
		err = row.Scan(
			&rem.Id,
			&rem.Text,
			&rem.UserID,
			&rem.RemindAt,
			&rem.RemindAtOriginal,
			&rem.Snoozed,
		)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, rem)

	}

	return reminders, nil
}

func (m *Manager) GetExpired(snoozePeriod string) ([]Reminder, error) {
	var err error
	var reminders []Reminder

	// todo: fix psql err about `invalid syntax near $1`
	row, err := m.db.Query(strings.Replace(getExpired, "%int%", snoozePeriod, 1))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var rem Reminder
		err = row.Scan(
			&rem.Id,
			&rem.Text,
			&rem.UserID,
			&rem.RemindAt,
			&rem.RemindAtOriginal,
			&rem.Snoozed,
		)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, rem)

	}

	return reminders, nil
}
