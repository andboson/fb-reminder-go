package main

import (
	"log"
	"time"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"
)

const checkInterval = 1 * time.Minute

func main() {
	var err error

	c, err := config()
	if err != nil {
		log.Fatalf("Err load config: %s", err)
	}

	db, err := InitDB(c)
	if err != nil {
		log.Fatalf("Err init db: %s", err)
	}

	rm := reminders.NewManager(db)
	fb := facebook.NewFBClient(c.FbToken)
	dfp := processor.NewDFProcessor("./config.json")

	scheduler := NewScheduler(rm, fb)
	scheduler.Start(c.SnoozePeriod, checkInterval)

	dispatcher := NewDispatcher(rm, fb, dfp)
	server := NewService(c.ServerAddress, dispatcher, c.XKey)

	if err = server.Serve(); err != nil {
		log.Fatalf("exit fatal")
	}
}
