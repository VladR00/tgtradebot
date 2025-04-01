package staffbot

import (
	"fmt"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
func StartMenuAdmin(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "admin panel")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("SupMenu", "Menu"),
			tgbotapi.NewInlineKeyboardButtonData("SupList", "SupList"),
			tgbotapi.NewInlineKeyboardButtonData("AddSup", "AddSup"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func AddSupButton(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	staff.AddSup = true
	staff.MapUpdateOrCreate()
	msg := tgbotapi.NewMessage(chatID, "Write ChatID of the future support")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "BackToMenuWithoutChanges"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func BackToMenuWithoutChanges(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	staff.AddSup = false
	staff.MapUpdateOrCreate()
	StartMenuAdmin(chatID, bot)
}

func SupListButton(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	suplist, err := database.OutputStaff()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error outputing suplist while read DB: %v", err), true)
	}
	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "SupList")

	var row []tgbotapi.InlineKeyboardButton

	for _, sup := range suplist {
		id := fmt.Sprintf("%d",sup.ChatID)
		if chatID == sup.ChatID{
			id = "You"
		}
		support := tgbotapi.NewInlineKeyboardButtonData(id, fmt.Sprintf("SupProfile%d",sup.ChatID))
		row = append(row, support)

		if len(row) == 2 {
			DefaultKeyboard = append(DefaultKeyboard, row)	
			row = []tgbotapi.InlineKeyboardButton{}
		}
		
	}

	if len(row) > 0 {
		DefaultKeyboard = append(DefaultKeyboard, row)
	}

	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "adminMenu"),
	}
	DefaultKeyboard = append(DefaultKeyboard, back)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(DefaultKeyboard...)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}