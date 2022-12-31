package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func telegramBotNotice(program, gitReversion, buildTime, buildGoVersion, flag, prefix string) {

	var localIp string
	ts := time.Now()
	bot, err := tgbotapi.NewBotAPI("5249320515:AAHibqLVtW69J6_OyJi1amDwXO1HfVTr3iw")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	format := "‼️‼️生产环境%s服务启动‼️‼️\r\n✅✅✅✅✅✅️\r\n⚠️Datetime: \t%s\r\n⚠️GitReversion: \t%s\r\n⚠️BuildTime: \t%s\r\n⚠️BuildGoVersion: \t%s\r\n⚠️Hostname: \t%s\r\n⚠️IP: \t%s\r\n⚠️Flag: \t%s\r\n⚠️Prefix: \t%s\n✨ ✨ ✨ ✨ ✨ ✨\r\n"
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
		return
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIp = ipnet.IP.String()
				break
			}
		}
	}

	msg := tgbotapi.NewMessage(-738052985, "")
	msg.Text = fmt.Sprintf(format, program, ts.Format("2006-01-02 15:04:05"), gitReversion, buildTime, buildGoVersion, hostname, localIp, flag, prefix)
	if _, err := bot.Send(msg); err != nil {
		fmt.Println("tgbot error : ", err.Error())
	}
}
