package main

import (
	//"sync"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartSupMenu(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "sup")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("menu", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}