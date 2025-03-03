package main

import (
	"log"
	//"strings"
	"fmt"

	database "tgbottrade/internal/database"
	help	 "tgbottrade/pkg/api/help"
	mainbot  "tgbottrade/internal/bot_main"
	supbot	 "tgbottrade/internal/bot_support/bot_user"
	staffbot "tgbottrade/internal/bot_support/bot_staff"
	//payment	 "tgbottrade/pkg/api/payment"
	config	 "tgbottrade/pkg/api/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"
)

var (
	bot                         *tgbotapi.BotAPI
	cryptoClient                *cryptobot.Client
)

func main() {
	config, err := config.LoadConfig("../config/config.json")
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
	bot.Debug = false 
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go mainbot.HandleMessageSwitchForMain(update, bot)
		}
		if update.CallbackQuery != nil {
			go mainbot.HandleCallBackSwitchForMain(update, bot, cryptoClient)
		}
	}
}
func supBotUpdates(bot *tgbotapi.BotAPI){
	bot.Debug = false 
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	
	for update := range updates {
		if update.Message != nil {
			staff, _ := database.ReadStaffByID(update.Message.Chat.ID)
			if (staff != nil){																		//Authorized
				if _, exists := database.StaffMap[update.Message.Chat.ID]; !exists{
					database.StaffMap[update.Message.Chat.ID] = *staff
				}
				fmt.Println(database.StaffMap)
				go staffbot.HandleMessageSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized
			} else {			
				user, _ := database.ReadUserByID(update.Message.Chat.ID)
				var err error
				if user == nil{		
					fmt.Println("user nil")																//Unauthorized
					if err := database.InsertNewUser(update.Message.Chat.ID, fmt.Sprintf("@%s",update.Message.Chat.UserName), update.Message.Chat.FirstName); err != nil{
						fmt.Println(err)
						continue;
					}
					if user, err = database.ReadUserByID(update.Message.Chat.ID); err != nil{
						help.NewMessage(update.Message.Chat.ID, bot, fmt.Sprintf("Error: %v\n Pleace contact us)))"), false)
						continue;
					}
				}		
				if _, exists := database.UserMap[update.Message.Chat.ID]; !exists{
					database.UserMap[update.Message.Chat.ID] = *user
				}		
				fmt.Println(database.UserMap)													//Unauthorized
				go supbot.HandleMessageSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
		if update.CallbackQuery != nil {
			staff, _ := database.ReadStaffByID(update.CallbackQuery.Message.Chat.ID)
			if (staff != nil){																			//Authorized
				if _, exists := database.StaffMap[update.CallbackQuery.Message.Chat.ID]; !exists{
					database.StaffMap[update.CallbackQuery.Message.Chat.ID] = *staff
					fmt.Println(database.StaffMap)
				}																
				fmt.Println(database.StaffMap)
				go staffbot.HandleCallBackSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized
			} else {
				var err error
				user, _ := database.ReadUserByID(update.CallbackQuery.Message.Chat.ID)
				if user == nil{		
					fmt.Println("user nil")																//Unauthorized
					if err := database.InsertNewUser(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("@%s",update.CallbackQuery.Message.Chat.UserName), update.CallbackQuery.Message.Chat.FirstName); err != nil{
						fmt.Println(err)
						continue;
					}
					if user, err = database.ReadUserByID(update.CallbackQuery.Message.Chat.ID); err != nil{
						help.NewMessage(update.CallbackQuery.Message.Chat.ID, bot, fmt.Sprintf("Error: %v\n Pleace contact us)))"), false)
						continue;
					}
				}
				if _, exists := database.UserMap[update.CallbackQuery.Message.Chat.ID]; !exists{
					database.UserMap[update.CallbackQuery.Message.Chat.ID] = *user
				}	
				fmt.Println(database.UserMap)
				go supbot.HandleCallBackSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
	}
}