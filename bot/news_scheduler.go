package bot

import (
	"bot-telegram/db"
	"bot-telegram/dtos"
	"bot-telegram/services/news"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// SendScheduledNews sends news to the users that have scheduled one for the given hour
func (nb *NewsBot) SendScheduledNews(hour int) error {
	usersInfo, err := nb.getScheduledUsers(hour)

	if err != nil {
		return err
	}

	if len(usersInfo) == 0 {
		return ErrNoUserScheduled
	}

	var barrier sync.WaitGroup
	for _, userInfo := range usersInfo {
		barrier.Add(1)
		go func(info dtos.UserInfo) {
			defer barrier.Done()
			nb.sendNewsForUser(info)

		}(userInfo)
	}

	barrier.Wait()

	return nil
}

func (nb *NewsBot) getScheduledUsers(hour int) ([]dtos.UserInfo, error) {
	var usersInfo []dtos.UserInfo
	scheduleData, err := nb.ScheduleDB.Get(hour)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, fmt.Errorf("%w: %v", errRetrievingScheduleInfo, err)
	}

	if errors.Is(err, db.ErrNotFound) {
		return usersInfo, nil
	}

	for userID, chatID := range scheduleData.UsersInfo {
		usersInfo = append(usersInfo, dtos.UserInfo{UserID: userID, ChatID: chatID})
	}

	return usersInfo, nil
}

// sendNewsForUser sends news to the given user
func (nb *NewsBot) sendNewsForUser(userInfo dtos.UserInfo) {
	userNewsData, err := nb.DB.Get(userInfo.UserID)
	if err != nil {
		fmt.Printf("unexpected error sending scheduled news for user %v: %v", userInfo, err)
		return
	}

	if len(userNewsData.WantedNews) == 0 {
		fmt.Printf("user %v dit not select any topic\n", userInfo.UserID)
		return
	}

	newsToSend, err := news.GetNew(strings.Join(userNewsData.WantedNews, ","))
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
		return
	}

	inlineKeyboard := nb.TelegramBot.CreateInlineKeyboard()
	message := news.GetSummarizedMessage(newsToSend, nb.GPTService)
	inlineKeyboard.AddURLButton("Go to the news", newsToSend.Url, 1)
	_, err = nb.TelegramBot.AdvancedMode().ASendMessage(userInfo.ChatID, message, "markdown", 0, false, false, nil, false, false, inlineKeyboard)
	if err != nil {
		panic(fmt.Sprintf("error sending scheduled news to user %v: %v", userInfo.UserID, err))
	}
}
