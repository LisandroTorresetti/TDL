package main

import (
	newsBot "bot-telegram/bot"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// The instance of the bot.

func main() {
	telegramNewsBot, err := newsBot.CreateNewsBot()
	if err != nil {
		fmt.Printf("error creating News Bot: %v", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("service error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Finish NewsBot successfully")

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	if err == nil {
		err = telegramNewsBot.Run()
		if err == nil {
			if err != nil {
				log.Printf("error while creating db %v", err)
			}
			telegramNewsBot.StartGoRoutines()
			telegramNewsBot.StartHandlers()
			fmt.Println("waiting for sigterm")
			<-c
			fmt.Println("exiting bot")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

}
