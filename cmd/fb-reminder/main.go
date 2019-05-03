package main

import (
	"log"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/reminders"

)

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
	fb := facebook.NewFBClient()

	server := NewService(c.ServerAddress, rm, fb)

	if err = server.Serve(); err != nil {
		log.Fatalf("exit fatal")
	}
}
