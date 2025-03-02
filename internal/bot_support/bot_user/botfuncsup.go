package supbot

import (
	"fmt"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upM := update.Message
	switch upM.Text {
		case "/initiate":
			db, err := database.OpenDB()
			if err != nil {
				fmt.Println(err)
				break; 
			}
			var isStaffExists bool 
			isStaffExists = true
			if !database.IsTableExists(db, "staff"){
				isStaffExists = false
			}
			db.Close()
			if !isStaffExists{
				if err = database.CreateTable("staff"); err != nil {
					help.NewMessage1(upM.Chat.ID, bot, fmt.Sprintf("%v", err), true)
					fmt.Println(err)
				}
				if err := database.InsertNewStaff(upM.Chat.ID, true, fmt.Sprintf("@%s",upM.Chat.UserName), upM.Chat.FirstName); err != nil{
					help.NewMessage1(upM.Chat.ID, bot, fmt.Sprintf("Error initiating: %v", err), true)
				}
				isStaffExists = true
			}
			StartMenu(upM.Chat.ID, bot)
		case "/start":
			StartMenu(upM.Chat.ID, bot)
	}
}

func HandleCallBackSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upCQ := update.CallbackQuery 
	switch upCQ.Data {
		case "Menu":
			StartMenu(upCQ.Message.Chat.ID, bot)
	}
}

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "пользователь")
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
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}