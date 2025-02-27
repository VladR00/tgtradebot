package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"fmt"

	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken string 	`json:"TokenTGbot"`
	TelegramSupBotToken string 	`json:"TokenSupbot"`
	CryptoBotToken   string 	`json:"TokenCryptobot"`
}

var (
	bot                         *tgbotapi.BotAPI
	cryptoClient                *cryptobot.Client
)

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cryptoClient = cryptobot.NewClient(cryptobot.Options{
		Testing:  true,
		APIToken: config.CryptoBotToken,
	})

	if err = CreateDB(); err != nil {
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
				if err := InsertNewUsersDB(upM.Chat.ID, fmt.Sprintf("@%s",upM.Chat.UserName), upM.Chat.FirstName); err != nil{
					fmt.Println(err)
				}
				StartMenu(upM.Chat.ID, bot)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			if strings.HasPrefix(upCQ.Data, "topup"){
				TopUp(bot, upCQ.Message.Chat.ID, cryptoClient, "TRX", strings.TrimPrefix(upCQ.Data, "topup"))
			}
			switch upCQ.Data {
			case "Menu":
				if err := InsertNewUsersDB(upCQ.Message.Chat.ID, fmt.Sprintf("@%s",upCQ.Message.Chat.UserName), upCQ.Message.Chat.FirstName); err != nil{
					fmt.Println(err)
				}
				StartMenu(upCQ.Message.Chat.ID, bot)
			case "Services":
				ServiceMenu(upCQ.Message.Chat.ID, bot)
			case "FAQ": 
				FAQ(upCQ.Message.Chat.ID, bot)
			case "Profile":
				Profile(upCQ.Message.Chat.ID, bot)
			case "нахуй":
				NewMessage(upCQ.Message.Chat.ID, bot, "нахуй", true)
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
				NewMessage(upM.Chat.ID, bot, "старт епт", false)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			switch upCQ.Data {
			case "Menu":
				NewMessage(upCQ.Message.Chat.ID, bot, "старт епт", false)
			}
		}
	}
}