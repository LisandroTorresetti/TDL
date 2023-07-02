package bot

import (
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/utils"
	"fmt"
	teleBot "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
	log "github.com/sirupsen/logrus"
)

var PossibleCategories = []string{"business", "entertainment", "environment", "food", "health", "politics", "science", "sports", "technology", "top", "tourism", "world"}

func StartHandlersOperations(deleteChan chan dtos.DeleteDataInformation, c, getWhitelistChan, getBlackListChan chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	go deleteData(deleteChan, database, bot)
	go getWhitelist(getWhitelistChan, database, bot)
	go getBlacklist(getBlackListChan, database, bot)
	go addCategory(c, database, bot)
}

func deleteData(c chan dtos.DeleteDataInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
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

func getBlacklist(c chan dtos.GetInformation, database db.DB[dtos.Data], bot *teleBot.Bot) {
	for {
		toRemove := <-c
		fmt.Printf("get blacklist was called with id: %d\n", toRemove)
		data, _ := database.Get(toRemove.Id)
		bot.SendMessage(toRemove.ToAnswer, fmt.Sprintf("your data blacklisted is: %v, press /start to see menu again", data.OmittedTopics), "", 0, false, false)
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
		go addCategoryForChat(toAddCategory, db, bot)
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
					for _, k := range s {
						if k != key {
							s = append(s, k)
						}
					}
					d.WantedNews = s
					db.Update(d)
					bot.SendMessage(info.ToAnswer, fmt.Sprintf("added %s to the categories wanted, current is: %+v", key, d.WantedNews), "", 0, false, false)
				}
			})
		}
	}
	bot.AdvancedMode().ASendMessage(info.ToAnswer, fmt.Sprintf("Select options below to remove categories, current selected are: %v", data.WantedNews), "", 0, false, false, nil, false, false, kb)
}
