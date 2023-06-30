package handlers

import (
	"fmt"

	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
)

func preferencesCallback(bot *bt.Bot) {
	CreatePreferencesHandler(bot)
}

func CreateStartHandler(bot *bt.Bot) {
    bot.AddHandler(`(/start|Back to menu)`, func(u *objs.Update) {
        
        kb := bot.CreateKeyboard(false,false,true,"Select one option ...")

        kb.AddButton("Give me a new", 1)
        kb.AddButton("Preferences", 1)

		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Welcome to Mokasin news bot!", "", u.Message.MessageId, false,false, nil, false, false, kb)

		if err != nil {
			fmt.Println(err)
		}
	},"private")
}