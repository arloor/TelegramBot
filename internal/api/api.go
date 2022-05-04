package api

import (
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
