package bot

import (
	"bot-telegram/db"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
)

type newsProvider interface {
	SummarizeNews(string) (string, error)
}

type NewsBot struct {
	TelegramBot *teleBot.Bot
	DB          db.DB[Data]
	GPTService  newsProvider
}

func (nb *NewsBot) Summarize(newsToSummarize string) (string, error) {
	return nb.GPTService.SummarizeNews(newsToSummarize)
}

func (nb *NewsBot) Run() error {
	err := nb.TelegramBot.Run()
	if err != nil {
		fmt.Printf("error initializating telegram bot: %v", err)
		return err
	}

	return nil
}
