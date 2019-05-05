package facebook

import (
	"fmt"

	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/andboson/fbbot"
	"golang.org/x/net/context"
)

type Manager interface {
	ShowMenu(ctx context.Context, userID string) error
	ShowCreateConfirm(ctx context.Context, userID string, rem reminders.Reminder) error
	ShowReminder(ctx context.Context, userID string, rem reminders.Reminder) error
	ShowForToday(ctx context.Context, userID string, rems []reminders.Reminder) error
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
	menu := fbbot.NewMenu()
	menu.AddMenuItems(
		fbbot.NewPostbackMenuItem("Create reminder", "reminder_create"),
		fbbot.NewPostbackMenuItem("Reminders for today", "show_today"),
		fbbot.NewPostbackMenuItem("Delete all reminders", "delete_all"),
	)

	return f.bot.AddPersistentMenus(menu)
}

func (f *FBClient) ShowMenu(ctx context.Context, userID string) (err error) {
	msg := fbbot.NewGenericMessage()
	msg.Bubbles = menuItems

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}

func (f *FBClient) ShowCreateConfirm(ctx context.Context, userID string, rem reminders.Reminder) (err error) {
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
	msg := fbbot.NewGenericMessage()
	btns := []fbbot.Button{
		{
			Type:    postbackType,
			Title:   "confirm",
			Payload: reminders.ActionString("confirm", rem),
		},
	}

	if !rem.Snoozed {
		btns = append(btns, fbbot.Button{
			Type:    postbackType,
			Title:   "snooze",
			Payload: reminders.ActionString("snooze", rem),
		})
	}

	msg.Bubbles = []fbbot.Bubble{
		{
			Title:    "Reminder alert",
			SubTitle: rem.Text,
			Buttons:  btns,
		},
	}

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}

func (f *FBClient) ShowForToday(ctx context.Context, userID string, rems []reminders.Reminder) (err error) {
	msg := fbbot.NewGenericMessage()
	msg.Bubbles = []fbbot.Bubble{}
	for _, rem := range rems {
		bubble := fbbot.Bubble{
			Title:    rem.Text,
			SubTitle: rem.RemindAtOriginal,
			Buttons: []fbbot.Button{
				{
					Type:    postbackType,
					Title:   "delete",
					Payload: reminders.ActionString("delete", rem),
				},
			},
		}

		msg.Bubbles = append(msg.Bubbles, bubble)
	}

	return f.bot.Send(fbbot.User{ID: userID}, msg)
}
