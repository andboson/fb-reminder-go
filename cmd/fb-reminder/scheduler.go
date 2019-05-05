package main

import (
	"context"
	"log"
	"time"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/reminders"
)

type Scheduler struct {
	rm            reminders.Reminderer
	fb            facebook.Manager
	snoozePeriod  string
	checkInterval time.Duration
	stopped       bool
}

func NewScheduler(rm reminders.Reminderer, fb facebook.Manager) *Scheduler {
	return &Scheduler{
		fb: fb,
		rm: rm,
	}
}

func (s *Scheduler) Start(snoozePeriod string, checkInterval time.Duration) {
	s.snoozePeriod = snoozePeriod
	s.checkInterval = checkInterval
	s.check()
}

func (s *Scheduler) Stop() {
	s.stopped = true
}

func (s *Scheduler) check() {
	ctx := context.Background()
	rems, err := s.rm.GetExpired(s.snoozePeriod)
	if err != nil {
		log.Printf("err run scheduler: %s", err)
		return
	}
	for _, rem := range rems {
		err = s.fb.ShowReminder(ctx, rem.UserID, rem)
		if err != nil {
			log.Printf("err send alert: %s   about: %+v", err, rem)
		}
	}

	if !s.stopped {
		time.AfterFunc(s.checkInterval, s.check)
	}
}
