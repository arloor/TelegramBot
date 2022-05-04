package main

import (
	"TelegramBot/internal/api"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var bot api.API

func init() {
	file := "/var/log/bot.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Lshortfile | log.Flags())
	bot = api.NewDefaultAPI()
}

func main() {
	log.Println(bot.GetMe())
	for {
		updates, err := bot.GetUpdates()
		if err == nil {
			for _, update := range updates {
				//解禁
				handleCallBackData(update)
				formatUpdate := api.NewFormatUpdate(&update)
				formatUpdate.Info()
				// 封禁
				if formatUpdate.NewChatMembers != nil {
					chatId := formatUpdate.Chat.ID
					for _, member := range formatUpdate.NewChatMembers {
						userId := member.ID
						if bot.SendWelcome(chatId, userId) == nil {
							log.Println("封禁用户", member.UserName)
							bot.RestricMember(chatId, userId, api.ChatPermissions{
								false, false, false, false, false, false, false, false,
							})
						}
					}
				}
			}
		}
	}
}

func handleCallBackData(update tgbotapi.Update) {
	callbackData := update.CallbackData()
	if callbackData != "" {
		split := strings.Split(callbackData, "@")
		if len(split) == 2 {
			userId, err := strconv.ParseInt(split[0], 10, 64)
			chatId, err := strconv.ParseInt(split[1], 10, 64)
			if err == nil && update.SentFrom().ID == userId {
				err := bot.RestricMember(chatId, userId, api.ChatPermissions{
					true, true, true, true, true, true, true, true,
				})
				if err == nil {
					bot.AnswerCallbackQuery(update.CallbackQuery.ID, "您可以发言了", false)
					return
				}
			}
		}
		bot.AnswerCallbackQuery(update.CallbackQuery.ID, "请不要瞎点", true)
	}

}
