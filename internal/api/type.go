package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

type FormatUpdate struct {
	User           *tgbotapi.User
	Chat           *tgbotapi.Chat
	Text           string
	NewChatMembers []tgbotapi.User
}

type AnswerCallbackQueryRequest struct {
	CallbackQueryId string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

func NewAnswerCallbackQueryRequest(callbackId string, text string, showAlert bool) AnswerCallbackQueryRequest {
	return AnswerCallbackQueryRequest{
		callbackId, text, showAlert,
	}

}

type ChatIdMsgId struct {
	ChatId           int64
	MsgId            int
	SendTimeInSecond int
}

type SendMessageRequest struct {
	ChatId      int64                         `json:"chat_id"`
	Text        string                        `json:"text"`
	ParseMode   string                        `json:"parse_mode"`
	ReplyMarkup tgbotapi.InlineKeyboardMarkup `json:"reply_markup"`
}

type RestrictMemberRequest struct {
	ChatId          int64 `json:"chat_id"`
	UserId          int64 `json:"user_id"`
	ChatPermissions `json:"permissions"`
	UtilDate        int64 `json:"until_date"`
}

func NewRestricMemberRequest(chatId int64, userId int64, permissions ChatPermissions) RestrictMemberRequest {
	return RestrictMemberRequest{
		chatId,
		userId,
		permissions,
		0,
	}
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

func NewWelcomeMessage(chatId int64, userId int64, userAlias string) SendMessageRequest {
	userAlias = strings.ReplaceAll(userAlias, "*", "")
	userAlias = strings.ReplaceAll(userAlias, "_", "")
	userAlias = strings.ReplaceAll(userAlias, "[", "")
	userAlias = strings.ReplaceAll(userAlias, "]", "")
	userAlias = strings.ReplaceAll(userAlias, "(", "")
	userAlias = strings.ReplaceAll(userAlias, ")", "")
	userAlias = strings.ReplaceAll(userAlias, "`", "")
	return SendMessageRequest{
		ParseMode: "Markdown",
		ChatId:    chatId,
		Text:      "欢迎[" + userAlias + "](tg://user?id=" + strconv.FormatInt(userId, 10) + ")来到本群组！\n您有120秒的时间点击*我不是机器人*以获取发言权限\n超期将被限制发言，*退群后重新加群*可再次看到该验证信息",
		ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: newWelcomeInlineKeyboard(chatId, userId),
		},
	}
}

func newWelcomeInlineKeyboard(chatId int64, userId int64) [][]tgbotapi.InlineKeyboardButton {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 1)
	row := make([]tgbotapi.InlineKeyboardButton, 1)
	userIdStr := strconv.FormatInt(userId, 10)
	chatIdStr := strconv.FormatInt(chatId, 10)
	data := userIdStr + "@" + chatIdStr
	row[0] = tgbotapi.InlineKeyboardButton{
		Text:         "我不是机器人",
		CallbackData: &data,
	}
	buttons[0] = row
	return buttons
}

func NewFormatUpdate(update *tgbotapi.Update) FormatUpdate {
	if update.ChatMember != nil { // 是成员更新
		return FormatUpdate{
			User:           update.ChatMember.NewChatMember.User,
			Chat:           &update.ChatMember.Chat,
			Text:           getText(update),
			NewChatMembers: getNewChatMembers(update),
		}
	} else { // 非成员更新
		return FormatUpdate{
			User:           update.SentFrom(),
			Chat:           update.FromChat(),
			Text:           getText(update),
			NewChatMembers: nil,
		}
	}
}
func (this FormatUpdate) Info() {
	if this.User != nil && this.Chat != nil && this.Text != "" {
		userAlias := BuildUserAlias(*this.User)
		var chatTitle string
		if this.Chat.Title != "" {
			chatTitle = this.Chat.Title
		} else {
			chatTitle = BuildUserAliasFromName(this.Chat.FirstName, this.Chat.LastName)
		}
		log.Printf("%s[%s]在【%s】说：%s\n", userAlias, this.User.UserName, chatTitle, this.Text)
	}
}

func getText(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.Text
	} else if update.EditedMessage != nil {
		return update.EditedMessage.Text
	} else if update.ChannelPost != nil {
		return update.ChannelPost.Text
	} else if update.EditedChannelPost != nil {
		return update.EditedChannelPost.Text
	}
	return ""
}

func getNewChatMembers(update *tgbotapi.Update) []tgbotapi.User {
	if update.Message != nil && update.Message.NewChatMembers != nil {
		return update.Message.NewChatMembers
	} else if update.ChatMember != nil && update.ChatMember.Chat.Type == "supergroup" {
		chatTitle := update.ChatMember.Chat.Title
		newChatMember := update.ChatMember.NewChatMember
		oldChatMember := update.ChatMember.OldChatMember
		userAlias := BuildUserAlias(*newChatMember.User)
		// 处理非restricted
		isNew2Chat := (oldChatMember.Status == "left" || oldChatMember.Status == "kicked") && (newChatMember.Status != "left" && newChatMember.Status != "kicked")
		// 处理restricted
		isNew2Chat = isNew2Chat || (!oldChatMember.IsMember && newChatMember.IsMember && oldChatMember.Status == "restricted" && newChatMember.Status == "restricted")
		if isNew2Chat {
			log.Printf("群组[%s]的%s[%s]从[%s,isMember:%t]变为[%s,isMember:%t]", chatTitle, userAlias, newChatMember.User.UserName, oldChatMember.Status, oldChatMember.IsMember, newChatMember.Status, newChatMember.IsMember)
			users := make([]tgbotapi.User, 1)
			users[0] = *newChatMember.User
			return users
		}
	}
	return nil
}
