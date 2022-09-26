package main

import (
	"TelegramBot/internal/api"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var bot api.API

func init() {
	file := "/var/log/bot.log"
	rollingFile := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50,
		MaxAge:     14,
		MaxBackups: 10,
		Compress:   false,
	}
	mw := io.MultiWriter(os.Stdout, rollingFile)
	log.SetOutput(mw)
	log.SetFlags(log.Lshortfile | log.Flags())
	bot = api.NewDefaultAPI()
}

func main() {
	toDeleteMsgs := make(chan api.ChatIdMsgId, 100)
	go func() {
		for {
			toDeleteMsg := <-toDeleteMsgs
			nowSecond := time.Now().Unix()
			var secondsToSleep int64 = 120 - nowSecond + int64(toDeleteMsg.SendTimeInSecond)
			if secondsToSleep >= 0 {
				time.Sleep(time.Duration(secondsToSleep) * time.Second)
			}
			err := bot.DeleteMessage(strconv.FormatInt(toDeleteMsg.ChatId, 10), toDeleteMsg.MsgId)
			if err != nil {
				log.Println("error clear msg", err)

			}

		}
	}()
	botInfo, err := bot.GetMe()
	if err != nil {
		log.Fatalln("登陆机器人失败，请检查网络和token", err)
	}
	log.Println("机器人登陆成功：", botInfo)
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
						chatIdMsgId, err := bot.SendWelcome(chatId, userId, userAlias)
						if err == nil {
							log.Printf("封禁用户 %d [%s][%s]", userId, member.UserName, userAlias)
							bot.RestricMember(chatId, userId, api.ChatPermissions{
								false, false, false, false, false, false, false, false,
							})
							if chatIdMsgId != nil {
								toDeleteMsgs <- *chatIdMsgId
							}
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
		if err := bot.DeleteMessage(strconv.FormatInt(update.FromChat().ID, 10), update.Message.MessageID); err == nil {
			bot.SendMessage(strconv.FormatInt(update.FromChat().ID, 10), "本群组不允许以*频道身份*发送消息！已删除此类消息！")
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
