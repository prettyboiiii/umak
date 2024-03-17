package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/prettyboiiii/umak/kamu"
	"github.com/robfig/cron/v3"
)

var (
	DIARY_NUMBER string
	CRON_TABS    []string
)

func init() {
	DIARY_NUMBER = os.Getenv("DIARY_NUMBER")
	if DIARY_NUMBER == "" {
		log.Fatal("DIARY_NUMBER is missing")
	}
	crontabs := os.Getenv("CRON_TABS")
	CRON_TABS = strings.Split(crontabs, ",")
	if len(CRON_TABS) == 0 {
		log.Fatal("CRON_TABS are missing")
	}
}

func main() {
	log.Println("Welcome to Umak!")
	kamuObj := kamu.New()
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.Default())))
	for _, crontab := range CRON_TABS {
		if _, err := c.AddFunc(crontab, func() {
			if err := kamuObj.GetPlaceInQueue(DIARY_NUMBER); errors.Is(err, kamu.SeesionEndedErr) {
				kamuObj.StartConversation()
				kamuObj.GetPlaceInQueue(DIARY_NUMBER)
			} else if err != nil {
				log.Fatal(err)
			}
		}); err != nil {
			log.Fatal(err)
		}
	}

	c.Start()

	// Keep the main function running
	// until the application is terminated
	select {}
}

type logger struct {
	log.Logger
}
