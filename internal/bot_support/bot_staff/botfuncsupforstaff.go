package staffbot

import (
	"fmt"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForAuthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upM := update.Message
	fmt.Printf("Handle message on support bot from staff: %s. From user: %s\n", upM.Text, upM.Chat.UserName)
	switch upM.Text {
		case "/start":
			StartMenu(upM.Chat.ID, bot)
	}
}

func HandleCallBackSwitchForAuthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upCQ := update.CallbackQuery
	fmt.Printf("Handle callback on support bot from staff: %s. From user: %s\n", upCQ.Data, upCQ.Message.Chat.UserName)
	switch upCQ.Data {
		case "Menu":
			StartMenu(upCQ.Message.Chat.ID, bot)
	}
}	

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "сапорт")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Menu", "Menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}