package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var supervisorToken = parseSupervisorToken()

func parseSupervisorToken() string {
	token := os.Getenv("SUPERVISOR_TOKEN")
	if token != "" {
		return token
	}

	log.Fatal("SUPERVISOR_TOKEN is missing")
	return ""
}

type EventPayload struct {
	PrayerId   int    `json:"prayer_id"`
	PrayerName string `json:"prayer_name"`
	IsReminder bool   `json:"is_reminder"`
}

func emitEvent(payload EventPayload) error {
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "http://supervisor/core/api/events/prayertime", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+supervisorToken)
	request.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("supervisor: Did not receive 200, instead received %d", response.StatusCode)
	}

	return nil
}
