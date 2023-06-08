package main

import (
	"bot-telegram/db"
	"fmt"
	"log"
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

type Data struct {
	Pepe int `json:"pepe"`
	Id   int `json:"id"`
}

type DeleteDataInformation struct {
	id       int
	toAnswer int
}

func (d Data) GetPrimaryKey() int {
	return d.Id
}

func deleteData(c chan DeleteDataInformation, database db.DB[Data]) {
	for {
		toRemove := <-c
		fmt.Printf("delete was called with id: %d\n", toRemove)
		database.Delete(toRemove.id)
		bot.SendMessage(toRemove.toAnswer, "your data was deleted, to start using our bot again send /start", "", 0, false, false)
	}
}

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
	deleteChan := make(chan DeleteDataInformation, 100)
	if err == nil {
		err = bot.Run()
		if err == nil {
			d, err := db.CreateDB[Data]("postgres", "pepe")
			if err != nil {
				log.Printf("error while creating db %v", err)
			}
			fmt.Println("entering")
			start(d, deleteChan)
			go deleteData(deleteChan, d)
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
func start(d db.DB[Data], dc chan DeleteDataInformation) {
	bot.AddHandler("/hi", func(u *objs.Update) {
		fmt.Println("hi was clicked")
		bot.SendMessage(u.Message.Chat.Id, "hello", "", u.Message.MessageId, false, false)
	}, "private")

	bot.AddHandler("/start", func(u *objs.Update) {
		kb := bot.CreateInlineKeyboard()
		kb.AddURLButton("click me to go to google", "google.com", 1)
		kb.AddCallbackButtonHandler("click me to remove all your data", "/hi", 2, func(update *objs.Update) {
			fmt.Println("delete was clicked")
			toRemove := DeleteDataInformation{
				id:       u.Message.From.Id,
				toAnswer: u.Message.Chat.Id,
			}
			bot.SendMessage(u.Message.Chat.Id, "your data is being removed", "", u.Message.MessageId, false, false)
			dc <- toRemove
		})
		di := Data{
			Pepe: 23,
			Id:   u.Message.From.Id,
		}

		d.Insert(di)

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "hi to you too 4", "", u.Message.MessageId, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Printf("error happened, %v\n", err)
		}
	}, "private")
}
