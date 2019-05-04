package facebook

import "github.com/andboson/fbbot"

const postbackType = "postback"

var menuItems = []fbbot.Bubble{
	{
		Title: "  Reminders menu",
		Buttons: []fbbot.Button{
			{
				Type:    postbackType,
				Title:   "Create reminder",
				Payload: "reminder_create",
			},
			{
				Type:    postbackType,
				Title:   "Reminders for today",
				Payload: "show_today",
			},
			{
				Type:    postbackType,
				Title:   "Delete all reminders",
				Payload: "delete_all",
			},
		},
	},
}
