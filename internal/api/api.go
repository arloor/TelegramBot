package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
)

type API struct {
	HttpClient http.Client
}

func NewDefaultAPI() API {
	return API{
		http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
}

func BuildUserAlias(user tgbotapi.User) string {
	userAlias := user.FirstName
	if user.LastName != "" {
		if userAlias != "" {
			userAlias += " "
		}
		userAlias += user.LastName
	}
	return userAlias
}
