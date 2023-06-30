package handlers

import (
	"fmt"

	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
)


func CreatePreferencesHandler(bot *bt.Bot) {
	bot.AddHandler("(/preferences|Preferences)", func(u *objs.Update) {
		kb := bot.CreateKeyboard(false,false,true,"Select one option ...")

        kb.AddButton("Topics", 1)
        kb.AddButton("Periodicity", 1)
		kb.AddButton("Back to menu", 2)

		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Select an option", "", u.Message.MessageId, false,false, nil, false, false, kb)

		if err != nil {
			fmt.Println(err)
		}

		// kb.AddCallbackButtonHandler("Topics", "/topics", 1, func(update *objs.Update) {
			
		// })

		// kb.AddCallbackButtonHandler("Periodicity", "/periodicity", 2, func(update *objs.Update) {
			
		// })

		// kb.AddCallbackButtonHandler("Back", "/start", 3, func(update *objs.Update) {
			
		// })
		
		// _, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "Please select one of the options below.", "", u.Message.MessageId, false, false, nil, false, false, kb)
		// if err != nil {
		// 	fmt.Printf("An error happened, %v\n", err)
		// }
	}, "private")
}