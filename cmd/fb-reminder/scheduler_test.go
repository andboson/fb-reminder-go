package main

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type scnedulerSuite struct {
	rm  reminders.Reminderer
	fb  facebook.Manager
	dfp processor.Processor
	ml  sync.Mutex

	suite.Suite
}

func (s *scnedulerSuite) Test_Scheduler() {
	invoked := false
	var rem0 reminders.Reminder
	rem1 := reminders.Reminder{
		UserID: "fbuid",
		Text:   "sometext",
	}

	s.rm.(*ReminderManagerMock).On("GetExpired", "").
		Return([]reminders.Reminder{rem1}, nil)

	s.fb.(*FBClientMock).On("ShowReminder", context.Background(), "fbuid", rem1).
		Return(nil).RunFn = func(arguments mock.Arguments) {
		log.Printf("-== %+v", arguments)
		r := arguments.Get(2).(reminders.Reminder)
		s.ml.Lock()
		invoked = true
		rem0 = r
		s.ml.Unlock()
	}

	scheduler := NewScheduler(s.rm, s.fb)
	scheduler.Start("", 1*time.Second)

	time.Sleep(2 * time.Second)

	s.ml.Lock()
	s.True(invoked)
	s.Equal(rem0.UserID, rem1.UserID)
	s.ml.Unlock()
}

func (s *scnedulerSuite) SetupSuite() {
	s.fb = new(FBClientMock)
	s.rm = new(ReminderManagerMock)
}

func (s *scnedulerSuite) TearDownSuite() {
	//s.srv.Stop()
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(scnedulerSuite))
}
