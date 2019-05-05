package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Dispatcherer interface {
	Dispatch(wr dialogflow.WebhookRequest) (proto.Message, error)
}

type Dispatcher struct {
	rm  reminders.Reminderer
	fb  facebook.Manager
	dfp processor.Processor
}

func NewDispatcher(rm reminders.Reminderer, fb facebook.Manager, dfp processor.Processor) *Dispatcher {
	return &Dispatcher{
		fb:  fb,
		rm:  rm,
		dfp: dfp,
	}
}

func (d *Dispatcher) Dispatch(wr dialogflow.WebhookRequest) (proto.Message, error) {
	var resp proto.Message
	var ctx = context.Background()
	var err error

	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetDisplayName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetParameters())

	fbClientID := extractFBClientID(wr)
	switch wr.GetQueryResult().GetIntent().GetDisplayName() {
	case "menu":
		err = d.fb.ShowMenu(ctx, fbClientID)
	case "reminder_create_set_text":
		rem := extractReminderParams(wr.GetQueryResult().GetOutputContexts())
		rem.UserID = fbClientID
		err = d.fb.ShowCreateConfirm(ctx, fbClientID, rem)
	case "reminder_action":
		resp = d.dfp.ReminderAction(ctx, fbClientID, wr.GetQueryResult(), d.rm)
	default:
		resp = d.dfp.HandleDefault(ctx, fbClientID)
	}

	if err != nil {
		log.Printf("err dispatch intent: %s", err)
	}

	return resp, nil
}

func extractReminderParams(contexts []*dialogflow.Context) reminders.Reminder {
	var rem = reminders.Reminder{}
	var err error

	for _, ctx := range contexts {
		if strings.Contains(ctx.GetName(), "/timeset") {
			date := time.Now()
			params := ctx.GetParameters().GetFields()
			if params["date_time"].GetStringValue() != "today" {
				date, err = time.Parse(time.RFC3339, params["date_time"].GetStringValue())
				if err != nil {
					log.Printf("err parse date_time param: %s (%s)", err, params["date_time"].GetStringValue())
				}
			}
			remindAt, err := time.Parse(time.RFC3339, params["time"].GetStringValue())
			if err != nil {
				log.Printf("err parse time param: %s (%s)", err, params["time"].GetStringValue())
			}
			remindAt = time.Date(
				remindAt.Year(),
				remindAt.Month(),
				date.Day(),
				remindAt.Hour(),
				remindAt.Minute(),
				remindAt.Second(),
				0,
				remindAt.Location())

			rem.RemindAtOriginal = remindAt.String()
			rem.Text = params["any"].GetStringValue()
			rem.RemindAt = remindAt

			break
		}
	}

	return rem
}

func extractFBClientID(wr dialogflow.WebhookRequest) string {
	var fbID string
	var odir = wr.GetOriginalDetectIntentRequest()
	if data, ok := odir.GetPayload().GetFields()["data"]; ok {
		if sender, ok := data.GetStructValue().GetFields()["sender"]; ok {
			senderStruct := sender.GetStructValue()
			if id, ok := senderStruct.GetFields()["id"]; ok {
				fbID = id.GetStringValue()
			}
		}
	}

	return fbID

}
