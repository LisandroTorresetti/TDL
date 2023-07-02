package main

import (
	newsBot "bot-telegram/bot"
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/services/gpt"
	"bot-telegram/services/gpt/config"
	"bot-telegram/services/news"
	"fmt"
	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// The instance of the bot.
var bot *bt.Bot

func main() {
	telegramNewsBot, err := newsBot.CreateNewsBot()
	if err != nil {
		fmt.Printf("error creating News Bot: %v", err)
		os.Exit(1)
	}
	bot = telegramNewsBot.TelegramBot

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
	deleteChan := make(chan dtos.DeleteDataInformation, 100)
	wi := make(chan dtos.GetInformation, 100)
	bi := make(chan dtos.GetInformation, 100)
	ci := make(chan dtos.GetInformation, 100)
	if err == nil {
		err = bot.Run()
		if err == nil {
			d, err := db.CreateDB[dtos.Data]("postgres", "pepe")
			if err != nil {
				log.Printf("error while creating db %v", err)
			}
			start(d, deleteChan, wi, bi, ci)
			newsBot.StartHandlersOperations(deleteChan, ci, wi, bi, d, bot)
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
func start(d db.DB[dtos.Data], dc chan dtos.DeleteDataInformation, wi, bi, ci chan dtos.GetInformation) {
	c, err := config.LoadConfig()
	if err != nil {
		panic("pone el error que quieras aca licha")
	}
	g := gpt.NewGPT(c, &http.Client{})
	bot.AddHandler("/hi", func(u *objs.Update) {
		fmt.Println("hi was called")
		bot.SendMessage(u.Message.Chat.Id, "hello", "", u.Message.MessageId, false, false)
	}, "private")

	bot.AddHandler("/start", func(u *objs.Update) {
		kb := bot.CreateInlineKeyboard()
		kb.AddURLButton("Click me to go to google", "google.com", 1)
		kb.AddCallbackButtonHandler("click me to remove all your data", "/hi", 2, func(update *objs.Update) {
			fmt.Println("delete was clicked")
			toRemove := dtos.DeleteDataInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			bot.SendMessage(u.Message.Chat.Id, "your data is being removed", "", u.Message.MessageId, false, false)
			dc <- toRemove
		})
		kb.AddCallbackButtonHandler("Check what was whitelisted", "/hi1", 3, func(update *objs.Update) {
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			fmt.Println("delete2 was clicked")
			wi <- toRemove
		})
		kb.AddCallbackButtonHandler("Check what was blacklisted", "/hi2", 4, func(update *objs.Update) {
			fmt.Println("delete 3was clicked")
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			bi <- toRemove
		})
		kb.AddCallbackButtonHandler("Get a summarized technology new", "/summarize", 5, func(update *objs.Update) {
			n, err := news.GetNew("technology")
			if err != nil {
				bot.SendMessage(u.Message.Chat.Id, "An error has ocurred, please try again later.", "", 0, false, false)
			} else {
				message := news.GetSummarizedMessage(n, g)
				bot.SendMessage(u.Message.Chat.Id, message, "markdown", 0, false, false)
			}
		})
		kb.AddCallbackButtonHandler("Add categories", "/summsarize", 5, func(update *objs.Update) {
			toRemove := dtos.GetInformation{
				Id:       u.Message.From.Id,
				ToAnswer: u.Message.Chat.Id,
			}
			ci <- toRemove
		})
		di := dtos.Data{
			OmittedTopics: []string{},
			WantedNews:    []string{},
			Id:            u.Message.From.Id,
		}

		d.Insert(di)

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Please select one of the options below.", "", u.Message.MessageId, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Printf("error happened, %v\n", err)
		}
	}, "private")
}
