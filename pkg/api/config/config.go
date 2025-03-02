package config

import (
	"os"
	"encoding/json"
)

type Config struct {
	TelegramBotToken string 	`json:"TokenTGbot"`
	TelegramSupBotToken string 	`json:"TokenSupbot"`
	CryptoBotToken   string 	`json:"TokenCryptobot"`
}

func LoadConfig(filename string) (Config, error) {
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
