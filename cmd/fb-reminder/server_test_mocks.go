package main

import (
	"context"

	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// reminders
type ReminderManagerMock struct {
	mock.Mock
}

func (rm *ReminderManagerMock) Create(rem reminders.Reminder) error {
	args := rm.Called(rem)
	return args.Error(0)
}

func (rm *ReminderManagerMock) GetByID(id int) (*reminders.Reminder, error) {
	args := rm.Called(id)
	return args.Get(0).(*reminders.Reminder), args.Error(1)
}

func (rm *ReminderManagerMock) DeleteByID(id int) error {
	args := rm.Called(id)
	return args.Error(0)
}

func (rm *ReminderManagerMock) SetSnooze(id int) error {
	args := rm.Called(id)
	return args.Error(0)
}

func (rm *ReminderManagerMock) DeleteAllByUser(userID string) error {
	args := rm.Called(userID)
	return args.Error(0)
}

func (rm *ReminderManagerMock) GetTodayByUser(userID string) ([]reminders.Reminder, error) {
	args := rm.Called(userID)
	return args.Get(0).([]reminders.Reminder), args.Error(1)
}

func (rm *ReminderManagerMock) GetExpired(snoozePeriod string) ([]reminders.Reminder, error) {
	args := rm.Called(snoozePeriod)
	return args.Get(0).([]reminders.Reminder), args.Error(1)
}

// fb client
type FBClientMock struct {
	mock.Mock
}

func (fb *FBClientMock) ShowMenu(ctx context.Context, userID string) error {
	args := fb.Called(ctx, userID)
	return args.Error(0)
}

func (fb *FBClientMock) ShowCreateConfirm(ctx context.Context, userID string, rem reminders.Reminder) error {
	args := fb.Called(ctx, userID, rem)
	return args.Error(0)
}

func (fb *FBClientMock) ShowReminder(ctx context.Context, userID string, rem reminders.Reminder) error {
	args := fb.Called(ctx, userID, rem)
	return args.Error(0)
}

func (fb *FBClientMock) ShowForToday(ctx context.Context, userID string, rems []reminders.Reminder) error {
	args := fb.Called(ctx, userID, rems)
	return args.Error(0)
}

func (fb *FBClientMock) SetupPersistentMenu(ctx context.Context) error {
	args := fb.Called(ctx)
	return args.Error(0)
}

///dialogflow Processor
type DialogFlowMock struct {
	mock.Mock
}

func (dp *DialogFlowMock) HandleDefault(ctx context.Context, fbClientID string) proto.Message {
	args := dp.Called(ctx, fbClientID)
	return args.Get(0).(proto.Message)
}

func (dp *DialogFlowMock) SimpleMessage(text string) proto.Message {
	args := dp.Called(text)
	return args.Get(0).(proto.Message)
}

func (dp *DialogFlowMock) ReminderAction(ctx context.Context, fbClientID string, qr *dialogflowpb.QueryResult, rm reminders.Reminderer) proto.Message {
	args := dp.Called(ctx, fbClientID, qr, rm)
	return args.Get(0).(proto.Message)
}

// Dispatcher mock
type DispatcherMock struct {
	mock.Mock
}

func (dp *DispatcherMock) Dispatch(wr dialogflow.WebhookRequest) (proto.Message, error) {
	args := dp.Called(wr)
	return args.Get(0).(proto.Message), args.Error(1)
}

// const
const showMenuIntentRequest = `{"responseId":"d143c793-211e-4bd8-a08e-e85a9a67393a-bca4db85","queryResult":{"queryText":"menu","parameters":{},"allRequiredParamsPresent":true,"fulfillmentMessages":[{"text":{"text":[""]}}],"intent":{"name":"projects/reminder-2dcf5/agent/intents/34f3bd91-bd29-4ac0-b3e2-3876bb934619",
"displayName":"menu"},"intentDetectionConfidence":1,"languageCode":"en"},"originalDetectIntentRequest":{"payload":
{
"data": {
		  "sender": {
					  "id":"fbclient_id1"
					}
		}
}},"session":"projects/reminder-2dcf5/agent/sessions/aef469aa-a720-476f-206e-8e607456790d"}`
