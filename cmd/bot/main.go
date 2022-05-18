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
				deleteChannelPost(update)
				handleCallBackData(update)
				formatUpdate := api.NewFormatUpdate(&update)
				formatUpdate.Info()
				// 封禁
				if formatUpdate.NewChatMembers != nil {
					chatId := formatUpdate.Chat.ID
					for _, member := range formatUpdate.NewChatMembers {
						userId := member.ID
						userAlias := api.BuildUserAlias(member)
						if bot.SendWelcome(chatId, userId, userAlias) == nil {
							log.Printf("封禁用户 %d [%s][%s]", userId, member.UserName, userAlias)
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

func deleteChannelPost(update tgbotapi.Update) {
	if update.Message != nil && update.Message.SenderChat != nil && update.Message.SenderChat.Type == "channel" {
		log.Println("检测到有人用channel身份发送消息，自动删除")
		bot.DeleteMessage(strconv.FormatInt(update.FromChat().ID, 10), update.Message.MessageID)
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
					true, true, false, true, true, false, true, false,
				})
				if err == nil {
					log.Println("解禁用户", userId)
					bot.AnswerCallbackQuery(update.CallbackQuery.ID, "您获得了发言权限😄", false)
					return
				}
			}
		}
		bot.AnswerCallbackQuery(update.CallbackQuery.ID, "该验证并不针对你，或者Bot的权限不足，请不要瞎搞🤢", false)
	}

}
