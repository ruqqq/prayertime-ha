package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var reminderOffset = parseReminderOffset()

func main() {
	log.Printf("\tSUPERVISOR_TOKEN: %s", supervisorToken)
	log.Printf("\tREMINDER_OFFSET: %d", reminderOffset)
	log.Printf("\tLOCATION_CODE: %s", prayerTimeLocationCode)
	log.Println()

	refresherCron := cron.New(cron.WithSeconds())
	_, _ = refresherCron.AddFunc("0 0 5 * * *", func() { refreshScheduleWorker() })
	refresherCron.Start()
	for _, entry := range refresherCron.Entries() {
		log.Printf("\tCron Scheduled: %s\n", entry.Next)
	}

	go refreshScheduleWorker()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()
	<-done
	log.Println("Exiting...")
}

func parseReminderOffset() int {
	reminderOffsetStr := os.Getenv("REMINDER_OFFSET")
	if reminderOffsetStr == "" {
		reminderOffsetStr = "-5"
	}
	reminderOffset, err := strconv.Atoi(reminderOffsetStr)
	if err != nil {
		log.Fatalf("Error parsing REMINDER_OFFSET: %s", reminderOffsetStr)
	}

	return reminderOffset
}

var myCron *cron.Cron

func refreshSchedule(noCache ...bool) error {
	if myCron != nil {
		myCron.Stop()
	}

	myCron = cron.New(cron.WithSeconds())

	now := time.Now()
	nowString := now.Format("Mon Jan 2 03:04:05PM -0700")
	log.Print(nowString + ":")

	result, err := getPrayerTimes(time.Now(), noCache...)
	if err != nil {
		return err
	}

	for index, timeStr := range result.Times {
		now := time.Now()

		pTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Printf("=== ERR: %v\n", err)
			return err
		}

		loc, _ := time.LoadLocation("Local")
		pTime = pTime.In(loc)

		hour := pTime.Hour()
		min := pTime.Minute()

		if hour < now.Hour() {
			continue
		}

		if hour == now.Hour() && min < now.Minute() {
			continue
		}

		nMinsBefore := pTime.Add(time.Duration(reminderOffset) * time.Minute)

		reminderSchedule := fmt.Sprintf("0 %d %d %d %d *", nMinsBefore.Minute(), nMinsBefore.Hour(), now.Day(), now.Month())
		schedule := fmt.Sprintf("0 %d %d %d %d *", min, hour, now.Day(), now.Month())

		//ind := index
		prayerId := index
		prayerName := waktuFromIndex(prayerId)
		_, _ = myCron.AddFunc(reminderSchedule, func() {
			err := emitEvent(EventPayload{
				PrayerId:   prayerId,
				PrayerName: prayerName,
				IsReminder: true,
			})
			if err != nil {
				log.Printf("Error: %s", err)
			}
		})
		_, _ = myCron.AddFunc(schedule, func() {
			err := emitEvent(EventPayload{
				PrayerId:   prayerId,
				PrayerName: prayerName,
				IsReminder: false,
			})
			if err != nil {
				log.Printf("Error: %s", err)
			}
		})
	}

	myCron.Start()
	for _, entry := range myCron.Entries() {
		log.Printf("\tCron Scheduled: %s\n", entry.Next)
	}

	log.Printf("=== DONE")

	return nil
}

func refreshScheduleWorker() {
	success := false

	// 3 retries max
	for i := 0; i < 3; i++ {
		err := refreshSchedule()
		if err == nil {
			success = true
			break
		}

		if i < 2 {
			duration := time.Duration(5*(i+1)) * time.Second
			log.Printf("=== RETRYING IN %v\n", duration)
			time.Sleep(duration)
		}
	}

	if success == false {
		log.Fatal("=== FAILED TO RETRIEVE PRAYER TIMES, EXITING")
	}
}

func waktuFromIndex(index int) string {
	if index == 0 {
		return "Subuh"
	} else if index == 1 {
		return "Syuruk"
	} else if index == 2 {
		return "Zuhur"
	} else if index == 3 {
		return "Asar"
	} else if index == 4 {
		return "Maghrib"
	} else if index == 5 {
		return "Isyak"
	}

	return "Unknown"
}
