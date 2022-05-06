## 电报机器人

```shell
bot_token=你的bot的token

rm -rf /var/TelegramBot
cd /var
git clone https://github.com/arloor/TelegramBot.git
cd /var/TelegramBot
go mod tidy
go install TelegramBot/cmd/bot
cat > /lib/systemd/system/bot.service <<EOF
[Unit]
Description=forwardproxy-Http代理
After=network-online.target
Wants=network-online.target

[Service]
WorkingDirectory=/opt/bot
EnvironmentFile=/opt/bot/env
ExecStart=/root/go/bin/bot
LimitNOFILE=100000
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
mkdir /opt/bot
cat > /opt/bot/env <<EOF
bot_token=${bot_token}
EOF
service bot stop
systemctl enable bot
service bot start
tail -f /var/log/bot.log
```
