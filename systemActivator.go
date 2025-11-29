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

var botToken string 
var stringChatId 	 string 

func main() {
	err := godotenv.Load() 
	if err != nil {
		log.Println(err)
	}
	botToken = os.Getenv("TOKEN_ID")
	stringChatId = os.Getenv("CHAT_ID")
	chatId := setString(stringChatId)
	initBot(botToken, chatId)
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

func initBot(botToken string, chatId int64) {
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
		if update.Message != nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
	switch update.Message.Command() {
		case "start":
			string := startService()
		}

	}
}

func startService() {
}
