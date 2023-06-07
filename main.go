package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
	objs "github.com/SakoDroid/telego/objects"
)

const token string = "fillWithToken"

// The instance of the bot.
var bot *bt.Bot

func main() {
	up := cfg.DefaultUpdateConfigs()

	cf := cfg.BotConfigs{
		BotAPI: cfg.DefaultBotAPI,
		APIKey: token, UpdateConfigs: up,
		Webhook:        false,
		LogFileAddress: cfg.DefaultLogFile,
	}

	var err error

	//Creating the bot using the created configs
	bot, err = bt.NewBot(&cf)
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	if err == nil {
		err = bot.Run()
		if err == nil {
			start()
			<-c
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

}
func start() {
	bot.AddHandler("/start", func(u *objs.Update) {

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.SendMessage(u.Message.Chat.Id, "licha puto", "", u.Message.MessageId, false, false)
		if err != nil {
			fmt.Println(err)
		}

	}, "private", "group")
}
