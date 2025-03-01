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
	staffbot "tgbottrade/staffbot"
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

	if err = database.CreateTable("users"); err != nil {
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
				mainbot.StartMenu(upM.Chat, bot)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			if strings.HasPrefix(upCQ.Data, "topup"){
				payment.TopUp(bot, upCQ.Message.Chat.ID, cryptoClient, "TRX", strings.TrimPrefix(upCQ.Data, "topup"))
			}
			switch upCQ.Data {
				case "Menu":
					mainbot.StartMenu(upCQ.Message.Chat, bot)
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
			staff, _ := database.ReadStaffByID(update.Message.Chat.ID)
			if (staff != nil){
				go staffbot.HandleMessageSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized
			} else {
				go supbot.HandleMessageSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
		if update.CallbackQuery != nil {
			staff, _ := database.ReadStaffByID(update.CallbackQuery.Message.Chat.ID)
			if (staff != nil){
				go staffbot.HandleCallBackSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized
			} else {
				go supbot.HandleCallBackSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
	}
}