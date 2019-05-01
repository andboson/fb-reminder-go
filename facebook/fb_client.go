package facebook

import "github.com/andboson/fb-reminder-go/reminder"

type FBManager interface {
	ShowMenu(userID string) error
	ShowCreateConfirm(userID string, rem reminder.Reminder) error
	ShowReminder(userID string, rem reminder.Reminder) error
	ShowForToday(userID string) error
	SetupPersistentMenu() error
}

type FBClient struct {
}

func NewFBClient() *FBClient {

	return &FBClient{}
}

func (f *FBClient) SetupPersistentMenu() (err error) {

	return
}

func (f *FBClient) ShowMenu(userID string) (err error) {

	return
}

func (f *FBClient) ShowCreateConfirm(userID string, rem reminder.Reminder) (err error) {

	return
}

func (f *FBClient) ShowReminder(userID string, rem reminder.Reminder) (err error) {

	return
}

func (f *FBClient) ShowForToday(userID string) (err error) {

	return
}
