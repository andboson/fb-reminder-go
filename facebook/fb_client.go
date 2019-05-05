package facebook

import (
	"fmt"
	"log"

	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/andboson/fbbot"
	"golang.org/x/net/context"
)

type FBManager interface {
	ShowMenu(ctx context.Context, userID string) error
	ShowCreateConfirm(ctx context.Context, userID string, rem reminders.Reminder) error
	ShowReminder(ctx context.Context, userID string, rem reminders.Reminder) error
	ShowForToday(ctx context.Context, userID string) error
	SetupPersistentMenu(ctx context.Context) error
}

type FBClient struct {
	bot *fbbot.Bot
}

func NewFBClient(pageToken string) *FBClient {
	bot := fbbot.New(0, "", "", pageToken)

	return &FBClient{
		bot: bot,
	}
}

func (f *FBClient) SetupPersistentMenu(ctx context.Context) (err error) {

	return
}

func (f *FBClient) ShowMenu(ctx context.Context, userID string) (err error) {
	msg := fbbot.NewGenericMessage()
	msg.Bubbles = menuItems

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}

func (f *FBClient) ShowCreateConfirm(ctx context.Context, userID string, rem reminders.Reminder) (err error) {
	log.Printf("rem: %+v", rem.String())
	msg := fbbot.NewGenericMessage()
	msg.Bubbles = []fbbot.Bubble{
		{
			Title:    "Confirm save reminder",
			SubTitle: fmt.Sprintf("Text: %s \n Time: %s", rem.Text, rem.RemindAtOriginal),
			Buttons: []fbbot.Button{
				{
					Type:    postbackType,
					Title:   "save",
					Payload: reminders.ActionString("save", rem),
				},
				{
					Type:    postbackType,
					Title:   "cancel",
					Payload: "cancel",
				},
			},
		},
	}

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}

func (f *FBClient) ShowReminder(ctx context.Context, userID string, rem reminders.Reminder) (err error) {

	return
}

func (f *FBClient) ShowForToday(ctx context.Context, userID string) (err error) {

	return
}
