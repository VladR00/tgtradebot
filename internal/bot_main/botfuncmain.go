package mainbot

import (
	"fmt"
	"strings"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	payment	 "tgbottrade/pkg/api/payment"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"
)

func HandleMessageSwitchForMain(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upM := update.Message;
	fmt.Printf("Handle message on main bot: %s. From user: %s\n", upM.Text, upM.Chat.UserName)
	switch upM.Text {
		case "/start":
			StartMenu(upM.Chat, bot)
	} 
}

func HandleCallBackSwitchForMain(update tgbotapi.Update, bot *tgbotapi.BotAPI, cryptoClient *cryptobot.Client){
	upCQ := update.CallbackQuery;
	fmt.Printf("Handle callback on main bot: %s. From user: %s\n", upCQ.Data, upCQ.Message.Chat.UserName)
	if strings.HasPrefix(upCQ.Data, "topup"){
		payment.TopUp(bot, upCQ.Message.Chat.ID, cryptoClient, "TRX", strings.TrimPrefix(upCQ.Data, "topup"))
	}
	switch upCQ.Data {
		case "Menu":
			 StartMenu(upCQ.Message.Chat, bot)
		case "Services":
			ServiceMenu(upCQ.Message.Chat.ID, bot)
		case "Profile":
			Profile(upCQ.Message.Chat.ID, bot)
	}
}

func StartMenu(upM *tgbotapi.Chat, bot *tgbotapi.BotAPI){
	chatID := upM.ID
	if err := database.InsertNewUser(chatID, fmt.Sprintf("@%s",upM.UserName), upM.FirstName); err != nil{
		fmt.Println(err)
	}
	go help.ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "hello muchahos")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Profile", "Profile"),
				tgbotapi.NewInlineKeyboardButtonData("Services", "Services"),
				tgbotapi.NewInlineKeyboardButtonURL("Support", "https://t.me/qweasdasdasddsaasdbot"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)	
}

func ServiceMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages(chatID, bot)
	_, err := database.ReadUserByID(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v",err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Fix", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)
		return
	}
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
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)	
}
func Profile(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages(chatID, bot)
	user, err := database.ReadUserByID(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v",err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Fix", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)
		return
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ID: %d\nLinkname: %s\nUsername: %s\nBalance: %d\nRegistration Time: %s", user.ChatID, user.LinkName, user.UserName, user.Balance, user.Time))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Top Up", "Services"),
				tgbotapi.NewInlineKeyboardButtonData("Back", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete(sent.Chat.ID, sent.MessageID)
}