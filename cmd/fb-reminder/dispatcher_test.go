package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type dispSuite struct {
	disp Dispatcherer

	rm  reminders.Reminderer
	fb  facebook.Manager
	dfp processor.Processor

	suite.Suite
}

func (s *dispSuite) Test_MenuIntentRequest() {
	s.fb.(*FBClientMock).On("ShowMenu", context.TODO(), "fbclient_id1").Return(nil)
	s.dfp.(*DialogFlowMock).On("HandleDefault", context.TODO(), "fbclient_id1").
		Return(&dialogflowpb.Intent_Message_SimpleResponse{})

	wr := dialogflow.WebhookRequest{}
	reader := bytes.NewReader([]byte(showMenuIntentRequest))
	err := jsonpb.Unmarshal(reader, &wr)
	s.Require().NoError(err)

	_, err = s.disp.Dispatch(wr)
	s.NoError(err)
}

func (s *dispSuite) SetupSuite() {
	s.fb = new(FBClientMock)
	s.rm = new(ReminderManagerMock)
	s.dfp = new(DialogFlowMock)

	s.disp = NewDispatcher(s.rm, s.fb, s.dfp)
}

func TestDispTestSuite(t *testing.T) {
	suite.Run(t, new(dispSuite))
}
