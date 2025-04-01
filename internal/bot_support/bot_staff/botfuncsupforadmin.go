package staffbot

import (
	"fmt"
	"strings"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForAdmin(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upM := update.Message
	fmt.Printf("Handle message from admin: %s. From user: %s\n", upM.Text, staff.UserName)

	switch upM.Text {
		case "/start":
			StartMenuAdmin(upM.Chat.ID, bot)
	}
}

func HandleCallBackSwitchForAdmin(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upCQ := update.CallbackQuery
	fmt.Printf("Handle callback from admin: %s. From user: %s\n", upCQ.Data, staff.UserName)

	switch upCQ.Data {
		case "Menu":
			StartMenuAdmin(upCQ.Message.Chat.ID, bot)
			return
	}

	switch {
		case strings.HasPrefix(upCQ.Data, "Accept"):
			AcceptTicket(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Accept"))
		
		case strings.HasPrefix(upCQ.Data, "Close"):
			CloseTicket(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Close"))
	}
}	

func StartMenuAdmin(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "admin")
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