package processor

import (
	"context"
	"log"

	"github.com/andboson/fb-reminder-go/reminders"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/golang/protobuf/proto"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Processor interface {
	HandleDefault(ctx context.Context, fbClientID string) proto.Message
	ReminderAction(ctx context.Context, fbClientID string, qr *dialogflowpb.QueryResult, rm reminders.Reminderer) proto.Message
	SimpleMessage(text string) proto.Message
}

type DFProcessor struct {
	authJSONFilePath string
	sessionClient    *dialogflow.SessionsClient
}

func NewDFProcessor(authFile string) *DFProcessor {
	var dp = DFProcessor{
		authJSONFilePath: authFile,
	}
	ctx := context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
		log.Fatal("Error in auth with Dialogflow")
	}
	dp.sessionClient = sessionClient
	return &dp
}

func (dp *DFProcessor) HandleDefault(ctx context.Context, fbClientID string) proto.Message {
	resp := &dialogflowpb.WebhookResponse{
		FulfillmentText: "Sorry, i didnt understand you",
	}

	return resp
}

func (dp *DFProcessor) ReminderAction(ctx context.Context, fbClientID string, qr *dialogflowpb.QueryResult, rm reminders.Reminderer) proto.Message {
	var err error
	var ffText = "done"
	var action = reminders.ActionFromString(qr.GetQueryText())

	switch action.Type {
	case "save":
		err = rm.Create(action.Alert)
	case "confirm":
		fallthrough
	case "delete":
		err = rm.DeleteByID(action.Alert.Id)
	case "snooze":
		err = rm.SetSnooze(action.Alert.Id)
	default:
		log.Printf("unknown action: %s", action.Type)
		ffText = "something wrong"
	}

	if err != nil {
		log.Printf("err action act: %s   obj: %+v", err, action)
		ffText = "error, try later"
	}
	resp := &dialogflowpb.WebhookResponse{
		FulfillmentText: ffText,
	}

	return resp
}

func (dp *DFProcessor) SimpleMessage(text string) proto.Message {
	return &dialogflowpb.WebhookResponse{
		FulfillmentText: text,
	}
}
