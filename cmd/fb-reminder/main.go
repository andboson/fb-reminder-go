package main

import (
	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/reminders"

	log "github.com/sirupsen/logrus"
)

func main() {
	var err error

	c, err := config()
	if err != nil {
		log.WithError(err).Fatalf("Err load config")
	}

	db, err := InitDB(c)
	if err != nil {
		log.WithError(err).Fatalf("Err init db")
	}

	rm := reminders.NewManager(db)
	fb := facebook.NewFBClient()

	server := NewService(c.ServerAddress, rm, fb)

	if err = server.Serve(); err != nil {
		log.Fatalf("exit fatal")
	}
}
