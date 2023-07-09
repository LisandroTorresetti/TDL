package scheduler

import (
	"fmt"
	"time"
)

type newsBot interface {
	SendScheduledNews(hour int) error
}

type Scheduler interface {
	Run()
}

// NewsScheduler entity that sends news to each user that have scheduled one for a given hour
type NewsScheduler struct {
	bot          newsBot
	newsInterval time.Duration
}

func NewNewsScheduler(newsBot newsBot, newsInterval time.Duration) *NewsScheduler {
	return &NewsScheduler{
		bot:          newsBot,
		newsInterval: newsInterval,
	}
}

// Run initialize the NewsScheduler
func (ns *NewsScheduler) Run() {
	currentTime := time.Now()
	nextHour := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour()+1, 0, 0, 0, currentTime.Location())
	duration := nextHour.Sub(currentTime)

	fmt.Printf("Waiting %v until the next o'clock hour\n", duration)
	delayTicker := time.NewTicker(duration)
	newsTicker := time.NewTicker(ns.newsInterval)
	newsTicker.Stop()

	for {
		select {
		case <-delayTicker.C:
			fmt.Println("The wait is over")
			newsTicker.Reset(ns.newsInterval)
			delayTicker.Stop()
			err := ns.bot.SendScheduledNews(time.Now().Hour())
			if err != nil {
				fmt.Printf("error sending shceduled news: %v\n", err)
			}

		case <-newsTicker.C:
			err := ns.bot.SendScheduledNews(time.Now().Hour())
			if err != nil {
				fmt.Printf("error sending shceduled news: %v\n", err)
			}
		}
	}
}
