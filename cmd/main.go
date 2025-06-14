package main

import (
	"log"
	"time"
	"fmt"
	"strconv"

	database "tgbottrade/internal/database"
	mainbot  "tgbottrade/internal/bot_main"
	supbot	 "tgbottrade/internal/bot_support/bot_user"
	staffbot "tgbottrade/internal/bot_support/bot_staff"
	config	 "tgbottrade/pkg/api/config"
	help 	 "tgbottrade/pkg/api/help"

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

	if err := database.InitiateTables(); err != nil {
		log.Fatalf("%v",err)
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

	if err := database.InitiateMaps(); err != nil {
		if err.Error() == "staff removed"{
			for _, value := range database.UserMap {
				if value.CurrentTicket != 0{
					continue
				}
				help.NewMessage1(value.ChatID, botsup, "Rare message. If you are reading this, it means that the support you were talking to was banned before the ticket could be closed. We recommend that you recreate the ticket and send it to this ticket number so that we are on the topic and don't have to rewrite your question.\nWrite: /start", true)
				ticketid, _ := strconv.ParseInt(value.UserName, 10, 64)
				ticket := database.TicketMap[ticketid]
				ticket.Status = "Closed"
				ticket.ClosingTime = time.Now().Unix()
				ticket.Update()
				delete(database.TicketMap, ticketid)	
				value.MapDelete()
			}
		} else {
			log.Fatalf("%v",err)
		}
	}

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
			if (staff != nil){	
				go staffbot.HandleMessageSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized
				
			} else {			
				if value, _ := database.ReadUserByID(update.Message.Chat.ID); value == nil{		
					user := database.User{
						ChatID:			update.Message.Chat.ID,
						LinkName:		fmt.Sprintf("@%s",update.Message.Chat.UserName),
						UserName:		update.Message.Chat.FirstName,
						Balance:		0,
						Time:			time.Now().Unix(),
						CurrentTicket:	0,
					}															
					if err := user.InsertNew(); err != nil{
						fmt.Println(err)
						continue;
					}
				}															
				go supbot.HandleMessageSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
		if update.CallbackQuery != nil {
			staff, _ := database.ReadStaffByID(update.CallbackQuery.Message.Chat.ID)
			if (staff != nil){				
				go staffbot.HandleCallBackSwitchForAuthorizedInTableStaff(update, bot, staff)			//Authorized																
			} else {
				if value, _ := database.ReadUserByID(update.CallbackQuery.Message.Chat.ID); value == nil{		
					user := database.User{
						ChatID:			update.CallbackQuery.Message.Chat.ID,
						LinkName:		fmt.Sprintf("@%s",update.CallbackQuery.Message.Chat.UserName),
						UserName:		update.CallbackQuery.Message.Chat.FirstName,
						Balance:		0,
						Time:			time.Now().Unix(),
						CurrentTicket:	0,
					}															
					if err := user.InsertNew(); err != nil{
						fmt.Println(err)
						continue;
					}
				}
				go supbot.HandleCallBackSwitchForUnauthorizedInTableStaff(update, bot)					//Unauthorized
			}
		}
	}
}