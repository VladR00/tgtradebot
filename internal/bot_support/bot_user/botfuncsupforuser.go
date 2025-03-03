package supbot

import (
	"fmt"

	staffbot "tgbottrade/internal/bot_support/bot_staff"
	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upM := update.Message
	fmt.Printf("Handle message on support bot from UnAuthorized: %s. From user: %s\n", upM.Text, upM.Chat.UserName)
	switch upM.Text {
		case "/start":
			StartMenu(upM.Chat.ID, bot)
	}
}

func HandleCallBackSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upCQ := update.CallbackQuery 
	fmt.Printf("Handle callback on support bot from UnAuthorized: %s. From user: %s\n", upCQ.Data, upCQ.Message.Chat.UserName)
	switch upCQ.Data {
		case "Menu":
			StartMenu(upCQ.Message.Chat.ID, bot)
		case "initiate":
			db, err := database.OpenDB()
			if err != nil {
				fmt.Println(err)
				return 
			}
			if !database.IsTableExists(db, "staff"){
				initiate(update, bot)
			} else {
				StartMenu(upCQ.Message.Chat.ID, bot)
			}
			db.Close()
	}
}

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	var defaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "пользователь")

	defaultKeyboard = [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("menu", "Menu")},
	}

	initiaterow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("initiate", "initiate"),
	}

	var isStaffExists bool 
	isStaffExists = true
	db, err := database.OpenDB()
	if err != nil {
		fmt.Println(err)
		return 
	}
	if !database.IsTableExists(db, "staff"){
		isStaffExists = false
	}
	db.Close()
	if !isStaffExists{
		defaultKeyboard = append(defaultKeyboard, initiaterow)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(defaultKeyboard...)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func initiate(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upCQ := update.CallbackQuery.Message
	db, err := database.OpenDB()
	if err != nil {
		fmt.Println(err)
		return 
	}
	defer db.Close()
	if err = database.CreateTable("staff"); err != nil {
		help.NewMessage1(upCQ.Chat.ID, bot, fmt.Sprintf("%v", err), true)
		fmt.Println(err)
	}
	if err := database.InsertNewStaff(upCQ.Chat.ID, true, fmt.Sprintf("@%s",upCQ.Chat.UserName), upCQ.Chat.FirstName); err != nil{
		help.NewMessage1(upCQ.Chat.ID, bot, fmt.Sprintf("Error initiating: %v", err), false)
	}
	staffbot.StartMenu(upCQ.Chat.ID, bot)
}