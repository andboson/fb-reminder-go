package reminders

import (
	"encoding/json"
	"log"
)

type Action struct {
	Type  string   `json:"type"`
	Alert Reminder `json:"alert"`
}

func ActionString(t string, r Reminder) string {
	b, _ := json.Marshal(Action{
		Type:  t,
		Alert: r,
	})

	return string(b)
}

func ActionFromString(actionString string) *Action {
	var a = new(Action)
	if err := json.Unmarshal([]byte(actionString), a); err != nil {
		log.Printf("err unmarshall action: %s", err)
	}

	return a
}
