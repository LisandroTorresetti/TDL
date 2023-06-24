package main

import (
	"bot-telegram/db"
	"bot-telegram/services"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	bt "github.com/SakoDroid/telego"
	cfg "github.com/SakoDroid/telego/configs"
	objs "github.com/SakoDroid/telego/objects"
	env "github.com/joho/godotenv"
)

const botTokenEnv = "TELEGRAM_BOT_TOKEN"

// The instance of the bot.
var bot *bt.Bot

type Data struct {
	NewsWanted []string `json:"news_wanted"`
	BlackList  []string `json:"black_list"`
	Id         int      `json:"id"`
}

type DeleteDataInformation struct {
	id       int
	toAnswer int
}

type GetInformation struct {
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

func getWhitelist(c chan GetInformation, database db.DB[Data]) {
	for {
		toRemove := <-c
		fmt.Printf("get whitelist was called with id: %d\n", toRemove)
		data, _ := database.Get(toRemove.id)
		bot.SendMessage(toRemove.toAnswer, fmt.Sprintf("your data whitelisted is: %v, press /start to see menu again", data.NewsWanted), "", 0, false, false)
	}
}

func getBlacklist(c chan GetInformation, database db.DB[Data]) {
	for {
		toRemove := <-c
		fmt.Printf("get blacklist was called with id: %d\n", toRemove)
		data, _ := database.Get(toRemove.id)
		bot.SendMessage(toRemove.toAnswer, fmt.Sprintf("your data blacklisted is: %v, press /start to see menu again", data.BlackList), "", 0, false, false)
	}
}
func main() {

	env_err := env.Load()

	if env_err != nil {
		fmt.Println(env_err)
		os.Exit(1)
	}

	token := os.Getenv(botTokenEnv)

	updateConfiguration := cfg.DefaultUpdateConfigs()

	cf := cfg.BotConfigs{
		BotAPI: cfg.DefaultBotAPI,
		APIKey: token, UpdateConfigs: updateConfiguration,
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
	wi := make(chan GetInformation, 100)
	bi := make(chan GetInformation, 100)
	if err == nil {
		err = bot.Run()
		if err == nil {
			d, err := db.CreateDB[Data]("postgres", "pepe")
			if err != nil {
				log.Printf("error while creating db %v", err)
			}
			start(d, deleteChan, wi, bi)
			go deleteData(deleteChan, d)
			go getBlacklist(bi, d)
			go getWhitelist(wi, d)
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
func start(d db.DB[Data], dc chan DeleteDataInformation, wi, bi chan GetInformation) {
	bot.AddHandler("/hi", func(u *objs.Update) {
		fmt.Println("hi was called")
		bot.SendMessage(u.Message.Chat.Id, "hello", "", u.Message.MessageId, false, false)
	}, "private")

	bot.AddHandler("/start", func(u *objs.Update) {
		kb := bot.CreateInlineKeyboard()
		kb.AddURLButton("Click me to go to google", "google.com", 1)
		kb.AddCallbackButtonHandler("click me to remove all your data", "/hi", 2, func(update *objs.Update) {
			fmt.Println("delete was clicked")
			toRemove := DeleteDataInformation{
				id:       u.Message.From.Id,
				toAnswer: u.Message.Chat.Id,
			}
			bot.SendMessage(u.Message.Chat.Id, "your data is being removed", "", u.Message.MessageId, false, false)
			dc <- toRemove
		})
		kb.AddCallbackButtonHandler("Check what was whitelisted", "/hi1", 3, func(update *objs.Update) {
			toRemove := GetInformation{
				id:       u.Message.From.Id,
				toAnswer: u.Message.Chat.Id,
			}
			fmt.Println("delete2 was clicked")
			wi <- toRemove
		})
		kb.AddCallbackButtonHandler("Check what was blacklisted", "/hi2", 4, func(update *objs.Update) {
			fmt.Println("delete 3was clicked")
			toRemove := GetInformation{
				id:       u.Message.From.Id,
				toAnswer: u.Message.Chat.Id,
			}
			bi <- toRemove
		})
		kb.AddCallbackButtonHandler("Summarize a hardcoded and short new", "/summarize", 5, func(update *objs.Update) {
			fmt.Println("Summarize new was clicked")
			newBody := services.GetRandomNew()
			summarizedNew := services.Summarize(newBody)
			bot.SendMessage(u.Message.Chat.Id, summarizedNew, "", 0, false, false)
		})
		di := Data{
			BlackList:  []string{},
			NewsWanted: []string{},
			Id:         u.Message.From.Id,
		}

		d.Insert(di)

		//Sends the message to the chat that the message has been received from. The message will be a reply to the received message.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Please select one of the options below.", "", u.Message.MessageId, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Printf("error happened, %v\n", err)
		}
	}, "private")
}
