package main

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
var (
	messagesMutex   sync.Mutex
)

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "hello muchahos")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Profile", "Profile"),
				tgbotapi.NewInlineKeyboardButtonData("Services", "Services"),
				tgbotapi.NewInlineKeyboardButtonData("F.A.Q.", "FAQ"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}

func FAQ(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "FAQ")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("иди нахуй", "нахуй"),

			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}
func ServiceMenu(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "Services")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("1000", "topup1000"),
				tgbotapi.NewInlineKeyboardButtonData("500", "topup500"),
				tgbotapi.NewInlineKeyboardButtonData("200", "topup200"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}
