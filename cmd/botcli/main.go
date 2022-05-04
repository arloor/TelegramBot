package main

import (
	"TelegramBot/internal/api"
	"log"
	"os"
	"strings"
)

var bot api.API

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Flags())
	bot = api.NewDefaultAPI()
}
func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalln("Usage: botcli @someUser message")
	} else {
		log.Println(args)
		bot.SendMessage(args[1], strings.Join(args[2:], " "))
	}
}
