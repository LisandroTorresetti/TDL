package handlers

import (
	bt "github.com/SakoDroid/telego"
	objs "github.com/SakoDroid/telego/objects"
)


func CreateTopicsHandler(bot *bt.Bot) {
	bot.AddHandler("(/topics|Topics)", func(u *objs.Update) {
		bot.SendMessage(u.Message.Chat.Id, `<b>List of topics</b><a href="#start">inline URL</a>`, "HTML", 0, false, false)
	}, "private")
}