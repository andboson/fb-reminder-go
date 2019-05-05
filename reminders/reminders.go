package reminders

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Reminderer interface {
	Create(rem Reminder) error
	GetByID(id int) (*Reminder, error)
	GetTodayByUser(userID string) ([]Reminder, error)
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
