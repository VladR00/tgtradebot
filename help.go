package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMessage(chatID int64, bot *tgbotapi.BotAPI,message string, needDelete bool){
	sent, err := bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		log.Println("Error sending message: ", err)
	}
	if (needDelete){
		go AddToDelete(sent.Chat.ID, sent.MessageID)
	}
}

func AddToDelete(chatID int64, messageID int) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	file, err := os.OpenFile("./messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error delete file: ",err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d:%d\n", chatID, messageID))
}
func ClearMessages(chatID int64, bot *tgbotapi.BotAPI) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Println("Error opening the delete file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var newData []string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		cid, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Printf("Error parsing chat ID: %s", err)
			continue
		}
		mid, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Error parsing message ID: %s", err)
			continue
		}
		if cid == chatID {
			bot.Request(tgbotapi.DeleteMessageConfig{
				ChatID:    cid,
				MessageID: mid,
			})
		} else {
			newData = append(newData, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from file: %s", err)
		return
	}

	file, err = os.Create("./messages.txt")
	if err != nil {
		log.Printf("Error opening the delete file for writing: %s", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newData {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Printf("Error writing to file: %s", err)
			return
		}
	}
	writer.Flush()
}