package facebook

import (
	"github.com/andboson/fb-reminder-go/reminders"

	"golang.org/x/net/context"
	"github.com/andboson/fbbot"
)

type FBManager interface {
	ShowMenu(ctx context.Context, userID string) error
	ShowCreateConfirm(userID string, rem reminders.Reminder) error
	ShowReminder(userID string, rem reminders.Reminder) error
	ShowForToday(userID string) error
	SetupPersistentMenu() error
}

type FBClient struct {
	bot *fbbot.Bot
}

func NewFBClient(pageToken string) *FBClient {
	bot := fbbot.New(0, "", "", pageToken)

	return &FBClient{
		bot:bot,
	}
}

func (f *FBClient) SetupPersistentMenu() (err error) {

	return
}

func (f *FBClient) ShowMenu(ctx context.Context, userID string) (err error) {
	msg := fbbot.NewGenericMessage()
//	msg.Text = "  Reminder menu"
	msg.Bubbles = menuItems

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}

func (f *FBClient) ShowCreateConfirm(userID string, rem reminders.Reminder) (err error) {

	return
}

func (f *FBClient) ShowReminder(userID string, rem reminders.Reminder) (err error) {

	return
}

func (f *FBClient) ShowForToday(userID string) (err error) {

	return
}
