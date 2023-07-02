package bot

import (
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/services/news"
	"bot-telegram/utils"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	CategoriesWanted = "categoriesWanted"
	DeleteData       = "deleteData"
	RemoveCategories = "removeCategories"
	AddCategories    = "addCategories"
	GetNew           = "getNews"
)

var PossibleCategories = []string{"business", "entertainment", "environment", "food", "health", "politics", "science", "sports", "technology", "top", "tourism", "world"}

func StartHandlersOperations(database db.DB[dtos.Data], bot *teleBot.Bot, provider news.Provider) map[string]chan dtos.GetInformation {
	deleteChan := getNewChan()
	getWhitelistChan := getNewChan()
	removeCategoryChan := getNewChan()
	addCategoryChan := getNewChan()
	getNewsChan := getNewChan()
	go deleteData(deleteChan, database, bot)
	go getWhitelist(getWhitelistChan, database, bot)
	go removeCategory(removeCategoryChan, database, bot)
	go addCategory(addCategoryChan, database, bot)
	go getWantedNews(getNewsChan, database, bot, provider)

	// ToDo there HAS to be a better way to do this
	return map[string]chan dtos.GetInformation{
		DeleteData:       deleteChan,
		CategoriesWanted: getWhitelistChan,
		RemoveCategories: removeCategoryChan,
		AddCategories:    addCategoryChan,
		GetNew:           getNewsChan,
	}
}

func getNewChan() chan dtos.GetInformation {
	return make(chan dtos.GetInformation, 100)
}

func deleteData(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toRemove := <-c
		fmt.Printf("delete was called with id: %d\n", toRemove)
		database.Delete(toRemove.Id)
		bot.SendMessage(toRemove.ToAnswer, "your data was deleted, to start using our bot again send /start", "", 0, false, false)
	}
}

func getWhitelist(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toRemove := <-c
		fmt.Printf("get whitelist was called with id: %d\n", toRemove)
		data, _ := database.Get(toRemove.Id)
		bot.SendMessage(toRemove.ToAnswer, fmt.Sprintf("your data whitelisted is: %v, press /start to see menu again", data.WantedNews), "", 0, false, false)
	}
}

func addCategory(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAddCategory := <-c
		log.Infof("checking of adding category, %d", toAddCategory.Id)
		go addCategoryForChat(toAddCategory, db, bot)
		bot.SendMessage(toAddCategory.ToAnswer, "processing request", "", 0, true, false)
	}
}

func addCategoryForChat(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	data, err := db.Get(info.Id)
	if err != nil {
		log.Errorf("couldn't get info for user, err: %v", err)
	}
	kb := bot.CreateInlineKeyboard()
	for i, k := range PossibleCategories {
		if !utils.Contains(k, data.WantedNews) {
			kb.AddCallbackButtonHandler(k, k, i%4+1, func(update *objs.Update) {
				d, err := db.Get(info.Id)
				key := update.CallbackQuery.Data
				if err != nil {
					log.Errorf("couldn't get info for user, err: %v", err)
					return
				}
				if utils.Contains(key, d.WantedNews) {
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("you already added this category to your wishlist: %s", key), "", 0, false, false)
				} else if len(d.WantedNews) > 4 {
					bot.SendMessage(info.ToAnswer, "you have too many categories already", "", 0, false, false)
				} else {
					d.WantedNews = append(d.WantedNews, key)
					db.Update(d)
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("added %s to the categories wanted, current is: %+v", key, d.WantedNews), "", 0, false, false)
				}
			})
		}
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Select options below to add categories, current selected are: %v", data.WantedNews), "", 0, false, false, nil, false, false, kb)
}

func removeCategory(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toAddCategory := <-c
		log.Infof("checking of removing category, %d", toAddCategory.Id)
		go removeCategoryForChat(toAddCategory, db, bot)
		bot.SendMessage(toAddCategory.ToAnswer, "processing request", "", 0, true, false)
	}
}

func removeCategoryForChat(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot) {
	data, err := db.Get(info.Id)
	if err != nil {
		log.Errorf("couldn't get info for user, err: %v", err)
	}
	kb := bot.CreateInlineKeyboard()
	for i, k := range PossibleCategories {
		if utils.Contains(k, data.WantedNews) {
			kb.AddCallbackButtonHandler(k, k, i+1, func(update *objs.Update) {
				d, err := db.Get(info.Id)
				key := update.CallbackQuery.Data
				if err != nil {
					log.Errorf("couldn't get info for user, err: %v", err)
					return
				}
				if !utils.Contains(key, d.WantedNews) {
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("you already removed this category from your wishlist: %s, selected: %+v", key, d.WantedNews), "", 0, false, false)
				} else {
					s := make([]string, 0, len(d.WantedNews)-1)
					for _, k := range d.WantedNews {
						if k != key {
							s = append(s, k)
						}
					}
					d.WantedNews = s
					db.Update(d)
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("removed %s to the categories wanted, current is: %+v", key, d.WantedNews), "", 0, false, false)
				}
			})
		}
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Select options below to remove categories, current selected are: %v", data.WantedNews), "", 0, false, false, nil, false, false, kb)
}

func getWantedNews(c chan dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot, gptService news.Provider) {
	for {
		toReturnNews := <-c
		log.Infof("user asking for news, %d", toReturnNews.Id)
		go getWantedNewsForUser(toReturnNews, db, bot, gptService)
		bot.SendMessage(toReturnNews.ToAnswer, "processing request", "", 0, true, false)
	}
}

func getWantedNewsForUser(info dtos.GetInformation, db db.DB[dtos.Data], bot *teleBot.Bot, gptService news.Provider) {
	data, err := db.Get(info.Id)
	topicsWanted := data.WantedNews
	if err != nil || len(topicsWanted) == 0 {
		bot.SendMessage(info.ToAnswer, "Topic not selected, searching for no category", "", 0, false, false)
	}
	n, err := news.GetNew(strings.Join(topicsWanted, ","))
	if err != nil {
		bot.SendMessage(info.ToAnswer, "An error has ocurred, please try again later.", "", 0, false, false)
	} else {
		kb := bot.CreateInlineKeyboard()
		message := news.GetSummarizedMessage(n, gptService)
		kb.AddURLButton("Ir a la noticia", n.Url, 1)
		bot.AdvancedMode().ASendMessage(info.ToAnswer, message, "markdown", 0, false, false, nil, false, false, kb)
	}
}
