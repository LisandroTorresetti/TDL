package bot

import (
	"bot-telegram/db"
	"bot-telegram/services/gpt"
	gptConfigLoader "bot-telegram/services/gpt/config"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	telegramConfig "github.com/SakoDroid/telego/configs"
	env "github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

const (
	botTokenEnv = "TELEGRAM_BOT_TOKEN"
	timeout     = 1
)

// ToDo: I Dont like this, refactor
type Data struct {
	NewsWanted []string `json:"news_wanted"`
	BlackList  []string `json:"black_list"`
	Id         int      `json:"id"`
}

func (d Data) GetPrimaryKey() int {
	return d.Id
}

// CreateNewsBot returns a NewsProvider with all the services it requires initialized
func CreateNewsBot() (*NewsBot, error) {
	database, err := db.CreateDB[Data]("postgres", "pepe")
	if err != nil {
		return nil, fmt.Errorf("error creating DB: %v", err)
	}

	err = env.Load()
	if err != nil {
		fmt.Println("error loading environment variables")
		return nil, err
	}

	token := os.Getenv(botTokenEnv)
	if token == "" {
		return nil, fmt.Errorf("error bot token missing")
	}

	updateConfiguration := telegramConfig.DefaultUpdateConfigs()
	botConfig := telegramConfig.BotConfigs{
		BotAPI: telegramConfig.DefaultBotAPI,
		APIKey: token, UpdateConfigs: updateConfiguration,
		Webhook:        false,
		LogFileAddress: telegramConfig.DefaultLogFile,
	}

	bot, err := teleBot.NewBot(&botConfig)
	if err != nil {
		fmt.Println("error creating telegram bot: %v", err)
		return nil, err
	}

	gptConfig, err := gptConfigLoader.LoadConfig()
	if err != nil {
		fmt.Printf("error loading gpt config: %v", err)
		return nil, err
	}

	// ToDo: check if we can do it different
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	gptService := gpt.NewGPT(gptConfig, client)

	return &NewsBot{
		TelegramBot: bot,
		DB:          database,
		GPTService:  gptService,
	}, nil
}
