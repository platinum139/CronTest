package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	mainJobsLogFile = "log.txt"
	controlJobsLogFile = "control_log.txt"
)

func main() {
	c := cron.New()

	// every minute, example 10:15:00
	if _, err := c.AddFunc("0-59 * * * *", func() {
		log := fmt.Sprintf("Minute job: %s\n", time.Now())
		// fmt.Print(log)
		WriteLog(controlJobsLogFile, log)
	}); err != nil {
		panic(err)
	}

	// every hour, example 10:00:00
	if _, err := c.AddFunc("0 * * * *", func() {
		log := fmt.Sprintf("Hour job: %s\n", time.Now())
		fmt.Print(log)
		WriteLog(mainJobsLogFile, log)
	}); err != nil {
		panic(err)
	}

	// every day, example 08.09.2022 00:00:00
	if _, err := c.AddFunc("0 0 * * *", func() {
		log := fmt.Sprintf("Day job: %s\n", time.Now())
		fmt.Print(log)
		WriteLog(mainJobsLogFile, log)
	}); err != nil {
		panic(err)
	}

	// every week, on Monday 05.09.2022 00:00:00
	if _, err := c.AddFunc("0 0 * * 1", func() {
		log := fmt.Sprintf("Week job: %s\n", time.Now())
		fmt.Print(log)
		WriteLog(mainJobsLogFile, log)
	}); err != nil {
		panic(err)
	}

	// every month, example 01.09.2022 00:00:00
	if _, err := c.AddFunc("0 0 1 * *", func() {
		log := fmt.Sprintf("Month job: %s\n", time.Now())
		fmt.Print(log)
		WriteLog(mainJobsLogFile, log)
	}); err != nil {
		panic(err)
	}

	c.Start()

	// graceful shutdown
	calls := make(chan os.Signal, 1)
	signal.Notify(calls, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-calls
		c.Stop()
	}()

	wg.Wait()
}

func WriteLog(logFileName, log string) {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(log)
	if err != nil {
		fmt.Println("Failed to write log to file: ", err)
	}
}
