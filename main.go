package main

import (
	"log"
	"strings"
	"fmt"

	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"
	database "tgbottrade/database"
	help	 "tgbottrade/help"
	mainbot  "tgbottrade/mainbot"
	supbot	 "tgbottrade/supbot"
	payment	 "tgbottrade/payment"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot                         *tgbotapi.BotAPI
	cryptoClient                *cryptobot.Client
)

func main() {
	config, err := help.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cryptoClient = cryptobot.NewClient(cryptobot.Options{
		Testing:  true,
		APIToken: config.CryptoBotToken,
	})

	if err = database.CreateDBusers(); err != nil {
		fmt.Println(err)
		return
	}

	botmain, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	log.Printf("Authorized on account %s", botmain.Self.UserName)
	
	botsup, err := tgbotapi.NewBotAPI(config.TelegramSupBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	log.Printf("Authorized on account %s", botsup.Self.UserName)

	go supBotUpdates(botsup)
	go mainBotUpdates(botmain)
	select {} //it's like for infinity "oo"
}

func mainBotUpdates(bot *tgbotapi.BotAPI){
	bot.Debug = true 
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			upM := update.Message;
			switch upM.Text {
			case "/start":
				if err := database.InsertNewUsersDB(upM.Chat.ID, fmt.Sprintf("@%s",upM.Chat.UserName), upM.Chat.FirstName); err != nil{
					fmt.Println(err)
				}
				mainbot.StartMenu(upM.Chat.ID, bot)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			if strings.HasPrefix(upCQ.Data, "topup"){
				payment.TopUp(bot, upCQ.Message.Chat.ID, cryptoClient, "TRX", strings.TrimPrefix(upCQ.Data, "topup"))
			}
			switch upCQ.Data {
			case "Menu":
				if err := database.InsertNewUsersDB(upCQ.Message.Chat.ID, fmt.Sprintf("@%s",upCQ.Message.Chat.UserName), upCQ.Message.Chat.FirstName); err != nil{
					fmt.Println(err)
				}
				mainbot.StartMenu(upCQ.Message.Chat.ID, bot)
			case "Services":
				mainbot.ServiceMenu(upCQ.Message.Chat.ID, bot)
			case "Profile":
				mainbot.Profile(upCQ.Message.Chat.ID, bot)
			}
		}
	}
}
func supBotUpdates(bot *tgbotapi.BotAPI){
	bot.Debug = true 
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			upM := update.Message;
			switch upM.Text {
			case "/start":
				supbot.StartMenu(upM.Chat.ID, bot)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			switch upCQ.Data {
			case "Menu":
				supbot.StartMenu(upCQ.Message.Chat.ID, bot)
			}
		}
	}
}