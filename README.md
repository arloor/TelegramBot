## 电报机器人

```
cd /var
git clone https://github.com/arloor/TelegramBot.git
cd /var/TelegramBot
git pull
go mod tidy
go install TelegramBot/cmd/bot
service bot restart
```
