package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)



type osInfo struct {
	hostname string 
	date string 
}

var helpMessage = `
Uso del bot: 
	/start 	<servicio>
	/status <servicio>
	/stop 	<servicio>
`


var botToken 		 string 
var stringChatId 	 string 

func main() {
	err := godotenv.Load() 
	if err != nil {
		log.Println(err)
	}
	botToken = os.Getenv("TOKEN_ID")
	stringChatId = os.Getenv("CHAT_ID")
	chatId := setString(stringChatId)
	info := setInfo()
	initBot(botToken, chatId, info)
}

func setString(stringChatId string) int64 {
	chatIdParsed, err := strconv.ParseInt(stringChatId, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return chatIdParsed
}

func setInfo() osInfo{
	osInfo := osInfo{
		hostname: getHostname(),
		date: time.Now().Format("2006-01-02 15:04:05"),
	}
	return osInfo
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return hostname
}

func initBot(botToken string, chatId int64, info osInfo) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30 
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message.Chat.ID != chatId {
			continue
		}
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
	    service := update.Message.CommandArguments()
	
		if service == "" {
			msg := tgbotapi.NewMessage(chatId, helpMessage)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			continue
		} 

	    var out string
		switch update.Message.Command() {
		case "start":
			out = startService(service)
		case "stop":
			out = stopService(service)
		case "status":
			out = statusService(service)
		default:
			out = helpMessage
		}
		data := fmt.Sprintf("Respuesta de:\n%s\n%s\n%s", info.hostname, info.date, out)
		msg := tgbotapi.NewMessage(chatId, data)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

	}
}

func startService(service string) string {
	cmd := exec.Command("systemctl", "start", service)
	out, err := cmd.CombinedOutput()
	time.Sleep(10 * time.Second)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	if string(out) ==  "" {
		return "El servicio se inicio de manera correcta"
	}
	return string(out)
}

func stopService(service string) string {
	cmd := exec.Command("systemctl", "stop", service)
	out, err := cmd.CombinedOutput()
	time.Sleep(10 * time.Second)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	if string(out) == "" {
		return "El servicio se desactivo de manera exitosa"
	}
	return string(out)

}

func statusService(service string) string {
	cmd := exec.Command("systemctl", "status", service)
	out, err := cmd.CombinedOutput()
	time.Sleep(10 * time.Second)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return string(out)
}

