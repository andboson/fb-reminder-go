package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type serverSuite struct {
	srv     *Service
	address string

	disp Dispatcherer

	suite.Suite
}

func (s *serverSuite) Test_MenuIntentRequest() {
	s.disp.(*DispatcherMock).On("Dispatch", mock.Anything).
		Return(&dialogflowpb.Intent_Message_SimpleResponse{}, nil)

	resp := s.request(showMenuIntentRequest, "POST", "/webhook")
	s.NotNil(resp)
	s.Equal(200, resp.StatusCode)

}

func (s *serverSuite) SetupSuite() {
	c := Config{}
	c.ServerAddress = "localhost:3002"
	s.address = c.ServerAddress

	s.disp = new(DispatcherMock)
	s.srv = NewService(c.ServerAddress, s.disp)

	go func() {
		err := s.srv.Serve()
		s.Require().NoError(err)
	}()
	time.Sleep(100 * time.Microsecond)
}

func (s *serverSuite) TearDownSuite() {
	//s.srv.Stop()
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
