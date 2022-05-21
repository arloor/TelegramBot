package api

import (
	"bytes"
	"encoding/json"
	"errors"
	TgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var updatesOffset int = 0

var botToken string = os.Getenv("bot_token")
var urlBase = "https://api.telegram.org/bot" + botToken + "/"

// 命令列表
var GetMe string = urlBase + "getMe"
var GetUpdates string = urlBase + "getUpdates"
var SendMessage string = urlBase + "sendMessage"
var DeleteMessage string = urlBase + "deleteMessage"
var RestrictChatMember string = urlBase + "restrictChatMember"
var AnswerCallbackQuery string = urlBase + "answerCallbackQuery"

func (this API) GetMe() (TgBot.User, error) {
	if res, err := this.HttpClient.Get(GetMe); err == nil {
		if res.StatusCode == 200 {
			response := TgBot.APIResponse{}
			all, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("read content err", err)
			}
			json.Unmarshal(all, &response)
			if response.Ok {
				user := TgBot.User{}
				json.Unmarshal(response.Result, &user)
				return user, nil
			}
		}
	}
	return TgBot.User{}, Error{"fail"}
}

func (receiver API) GetUpdates() ([]TgBot.Update, error) {
	return receiver.GetUpdateWithOffset(0)

}

func (receiver API) GetUpdateWithOffset(offset int) ([]TgBot.Update, error) {
	if offset != 0 {
		updatesOffset = offset
	}
	var body string = "offset=" + strconv.Itoa(updatesOffset) + "&allowed_updates=[\"message\",\"edited_message\",\"channel_post\",\"edited_channel_post\",\"inline_query\",\"chosen_inline_result\",\"callback_query\",\"shipping_query\",\"pre_checkout_query\",\"poll\",\"poll_answer\",\"my_chat_member\",\"chat_member\",\"chat_join_request\"]"
	if res, err := receiver.HttpClient.Post(GetUpdates, "application/x-www-form-urlencoded", strings.NewReader(body)); err == nil {
		if res.StatusCode == 200 {
			response := TgBot.APIResponse{}
			all, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("read content err", err)
			}
			json.Unmarshal(all, &response)
			if response.Ok {
				updates := make([]TgBot.Update, 1)
				json.Unmarshal(response.Result, &updates)
				len := len(updates)
				if len != 0 {
					updatesOffset = updates[len-1].UpdateID + 1
				}
				return updates, nil
			}
		}
	}
	return nil, Error{"fail"}
}

func (this API) SendWelcome(chatId int64, userId int64, userAlias string) error {
	message := NewWelcomeMessage(chatId, userId, userAlias)
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//log.Println(string(body))
	res, err := this.HttpClient.Post(SendMessage, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		return Error{"发送欢迎信息失败"}
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("read content err")
		return err
	}
	response := TgBot.APIResponse{}
	err = json.Unmarshal(all, &response)
	if err != nil || !response.Ok {
		log.Println("解析响应失败", err, string(all))
		return err
	}
	return nil
}

func (this API) DeleteMessage(chatId string, messageId int) error {
	message := make(map[string]interface{}, 3)
	message["chat_id"] = chatId
	message["message_id"] = messageId
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	res, err := this.HttpClient.Post(DeleteMessage, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		return Error{"删除信息失败"}
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("read content err")
		return err
	}
	response := TgBot.APIResponse{}
	err = json.Unmarshal(all, &response)
	if err != nil || !response.Ok {
		log.Println("解析响应失败", err, string(all))
		return errors.New("删除消息失败")
	}
	return nil
}

func (this API) SendMessage(userName string, text string) error {
	message := make(map[string]string, 3)
	message["chat_id"] = userName
	message["text"] = text
	message["parse_mode"] = "Markdown"
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	res, err := this.HttpClient.Post(SendMessage, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		return Error{"发送信息失败"}
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("read content err")
		return err
	}
	response := TgBot.APIResponse{}
	err = json.Unmarshal(all, &response)
	if err != nil || !response.Ok {
		log.Println("解析响应失败", err, string(all))
		return err
	}
	return nil
}

func (this API) RestricMember(chatId int64, userId int64, permissions ChatPermissions) error {
	request := NewRestricMemberRequest(chatId, userId, permissions)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	post, err := this.HttpClient.Post(RestrictChatMember, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	all, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return err
	}
	response := TgBot.APIResponse{}
	err = json.Unmarshal(all, &response)
	if err != nil {
		return err
	}
	if !response.Ok {
		log.Println("调整用户权限失败", string(all))
		return Error{"error restrict！"}
	}
	return nil
}

func (this API) AnswerCallbackQuery(callbackId string, text string, showAlert bool) error {
	request := NewAnswerCallbackQueryRequest(callbackId, text, showAlert)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	post, err := this.HttpClient.Post(AnswerCallbackQuery, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	all, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return err
	}
	response := TgBot.APIResponse{}
	err = json.Unmarshal(all, &response)
	if err != nil {
		return err
	}
	if !response.Ok {
		log.Println("AnswerCallbackQuery error", string(all))
		return Error{"error AnswerCallbackQuery！"}
	}
	return nil
}
