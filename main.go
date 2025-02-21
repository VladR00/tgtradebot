package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

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
	Supbot						*tgbotapi.BotAPI
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

	bot, err = tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	Supbot, err = tgbotapi.NewBotAPI(config.TelegramSupBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	cryptoClient = cryptobot.NewClient(cryptobot.Options{
		Testing:  true,
		APIToken: config.CryptoBotToken,
	})

	log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true 
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			upM := update.Message;
			switch upM.Text {
			case "/start":
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
				StartMenu(upCQ.Message.Chat.ID, bot)
			case "Services":
				ServiceMenu(upCQ.Message.Chat.ID, bot)
			case "FAQ": 
				FAQ(upCQ.Message.Chat.ID, bot)
			case "нахуй":
				NewMessage(upCQ.Message.Chat.ID, bot, "нахуй", true)
			}
		}
	}
}