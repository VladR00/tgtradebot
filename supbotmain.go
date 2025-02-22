package main

import(
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
func supbotupdates(){
	log.Printf("Authorized on sup account %s", Supbot.Self.UserName)

	Supbot.Debug = true 
	us := tgbotapi.NewUpdate(0)
	us.Timeout = 60
	
	Supupdates := bot.GetUpdatesChan(us)

	for update := range Supupdates {
		if update.Message != nil {
			upM := update.Message;
			switch upM.Text {
			case "/start":
				SupStartMenu(upM.Chat.ID, Supbot)
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			switch upCQ.Data {
			case "Menu":
				SupStartMenu(upCQ.Message.Chat.ID, Supbot)
			}
		}
	}
}

func SupStartMenu(chatID int64, bot *tgbotapi.BotAPI){
	//go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "supmenu")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("1000", "topup1000"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		// sent, err := bot.Send(msg)
		// if err != nil {
		// 	log.Println("Error sending start menu: ", err)
		// }
		//go AddToDelete(sent.Chat.ID, sent.MessageID)	
}