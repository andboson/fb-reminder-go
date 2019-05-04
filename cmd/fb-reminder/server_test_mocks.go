package main

import (
	"context"
	"github.com/andboson/fb-reminder-go/reminders"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
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

func (rm *ReminderManagerMock) GetTodayByUser(userID string) ([]reminders.Reminder, error) {
	args := rm.Called(userID)
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

func (fb *FBClientMock) ShowCreateConfirm(userID string, rem reminders.Reminder) error {
	args := fb.Called(userID, rem)
	return args.Error(0)
}

func (fb *FBClientMock) ShowReminder(userID string, rem reminders.Reminder) error {
	args := fb.Called(userID, rem)
	return args.Error(0)
}

func (fb *FBClientMock) ShowForToday(userID string) error {
	args := fb.Called(userID)
	return args.Error(0)
}

func (fb *FBClientMock) SetupPersistentMenu() error {
	args := fb.Called()
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
