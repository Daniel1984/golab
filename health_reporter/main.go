package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var statusEndpoint = "https://hooks.slack.com/services/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"

var pl = []string{
	"https://api.hostmaker.co/health",
	"https://hostmaker.co",
	"https://accounts.hostmaker.co",
	"https://ops-frontend.hostmaker.co",
	"https://ops.hostmaker.co",
	"https://pricing.hostmaker.co",
	"https://stay.hostmaker.co",
	"https://onboarding.hostmaker.co",
}

var sl = []string{
	"https://api-staging.hostmaker.co/health",
	"https://staging.hostmaker.co",
	"https://accounts-staging.hostmaker.co",
	"https://ops-frontend-staging.hostmaker.co",
	"https://ops-staging.hostmaker.co",
	"https://pricing-staging.hostmaker.co",
	"https://stay-staging.hostmaker.co",
	"https://onboarding-staging.hostmaker.co",
}

var dl = []string{
	"https://api-dev.hostmaker.co/health",
	"https://dev.hostmaker.co",
	"https://accounts-dev.hostmaker.co",
	"https://ops-frontend-dev.hostmaker.co",
	"https://ops-dev.hostmaker.co/sign_in",
	"http://stay-dev.hostmaker.co",
	"http://onboarding-dev.hostmaker.co",
}

// SlackAttachment will contain attachments map
type SlackAttachment struct {
	Color string `json:"color"`
	Text  string `json:"text"`
}

// SlackMessages will contain attachments map
type SlackMessages struct {
	Attachments []SlackAttachment `json:"attachments"`
}

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch env := r.FormValue("env"); env {
	case "staging":
		checkAllAppsStatus(sl)
	case "dev":
		checkAllAppsStatus(dl)
	default:
		checkAllAppsStatus(pl)
	}
}

func checkAllAppsStatus(links []string) {
	lastCheck := len(links) - 1

	attachments := []SlackAttachment{}
	ch := make(chan SlackAttachment)

	for i, link := range links {
		go getAppStatus(link, ch)
	}

	for status := range ch {
		attachments = append(attachments, status)

		if i == lastCheck {
			messages := SlackMessages{Attachments: attachments}
			jsonVal, _ := json.Marshal(messages)
			http.Post(statusEndpoint, "application/json", bytes.NewBuffer(jsonVal))
		}
	}
}

func getAppStatus(link string, ch chan SlackAttachment) {
	start := time.Now()
	_, err := http.Get(link)
	timeElapsed := time.Since(start)

	if err != nil {
		ch <- SlackAttachment{
			Color: "#F35A00",
			Text:  fmt.Sprintf("%s - DOWN, Response time - %s", link, timeElapsed),
		}
	}

	ch <- SlackAttachment{
		Color: "#7CD197",
		Text:  fmt.Sprintf("%s - UP, Response time - %s", link, timeElapsed),
	}
}
