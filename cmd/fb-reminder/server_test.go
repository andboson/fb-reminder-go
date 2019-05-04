package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/stretchr/testify/suite"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type serverSuite struct {
	srv     *Service
	address string

	rm  reminders.Reminderer
	fb  facebook.FBManager
	dfp processor.Processor

	suite.Suite
}

func (s *serverSuite) Test_MenuIntentRequest() {
	s.fb.(*FBClientMock).On("ShowMenu", context.TODO(), "fbclient_id1").Return(nil)
	s.dfp.(*DialogFlowMock).On("HandleDefault", context.TODO(), "fbclient_id1").
		Return(&dialogflowpb.Intent_Message_SimpleResponse{})

	resp := s.request(showMenuIntentRequest, "POST", "/webhook")
	s.NotNil(resp)
	s.Equal(200, resp.StatusCode)

}

func (s *serverSuite) SetupSuite() {
	c := Config{}
	c.ServerAddress = "localhost:3002"
	s.address = c.ServerAddress

	s.fb = new(FBClientMock)
	s.rm = new(ReminderManagerMock)
	s.dfp = new(DialogFlowMock)
	s.srv = NewService(c.ServerAddress, s.rm, s.fb, s.dfp)

	go func() {
		err := s.srv.Serve()
		s.Require().NoError(err)
	}()
	time.Sleep(500 * time.Microsecond)
}

func (s *serverSuite) TearDownSuite() {
	s.srv.Stop()
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(serverSuite))
}

func (s *serverSuite) request(json, method, url string) *http.Response {
	r, err := http.NewRequest(method, fmt.Sprintf("http://%s%s", s.address, url), strings.NewReader(json))
	if err != nil {
		s.FailNow(err.Error(), "error create new request")
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		s.FailNow(err.Error(), "error do request")
	}

	return res
}
