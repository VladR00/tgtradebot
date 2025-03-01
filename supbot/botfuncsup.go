package supbot

import (
	"fmt"

	//database "tgbottrade/database"
	help 	 "tgbottrade/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages(chatID, bot)
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
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)	
}