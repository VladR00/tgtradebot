package supbot

import (
	"fmt"
	"time"
	"strings"

	staffbot "tgbottrade/internal/bot_support/bot_staff"
	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upM := update.Message
	fmt.Printf("Handle message on support bot from UnAuthorized: %s. From user: %s\n", upM.Text, upM.Chat.UserName)
	
	if value, exists := database.UserMap[upM.Chat.ID]; exists{
		if value.UserName == "............................................................................."{
			value.UserName = upM.Text
			value.CurrentTicket = 0
			value.MapUpdateOrCreate()
			CreateTicket(value, bot)
		} else if (value.CurrentTicket != 0){
			message := database.TicketMessage{
				TicketID:	value.CurrentTicket,
				Support:	0,
				ChatID:		value.ChatID,
				UserName:	value.UserName,
				MessageID:	upM.MessageID,
				Time:		time.Now().Unix(),		
			}
			if err := message.InsertNew(); err != nil{
				fmt.Println(err)
				help.NewMessage(value.ChatID, bot, fmt.Sprintf("%v", err),false)
				return
			}
			ticket, err := database.ReadTicketByID(value.CurrentTicket)
			if err != nil {
				fmt.Println(err)
				help.NewMessage(value.ChatID, bot, fmt.Sprintf("%v", err),false)
				return
			}
			if ticket.SupChatID == 0 && ticket.Status == "Notificate"{
				fmt.Println("notificate:")
				go staffbot.NotificateSups(value, bot)
			} else {
				msg := tgbotapi.NewMessage(ticket.SupChatID, fmt.Sprintf("User %s send message by ticket %d, prefered language: %s", ticket.UserName, ticket.TicketID, ticket.Language))
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Communicate", fmt.Sprintf("Move%d",ticket.TicketID)),
					),
				)
				msg.ReplyMarkup = keyboard
				sent, err := bot.Send(msg)
				if err != nil {
					fmt.Println("SupChatID:", ticket.SupChatID)

					fmt.Println("Error sending start menu: ", err)
				}
				go help.AddToDelete1(ticket.SupChatID, sent.MessageID)
			}
		}
	}
	switch upM.Text {
		case "/start":
			StartMenu(upM.Chat.ID, bot)
	}
}

func HandleCallBackSwitchForUnauthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upCQ := update.CallbackQuery 
	fmt.Printf("Handle callback on support bot from UnAuthorized: %s. From user: %s\n", upCQ.Data, upCQ.Message.Chat.UserName)
	switch {
		case strings.HasPrefix(upCQ.Data, "Language"):
			CreateTicketName(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Language"))
		case upCQ.Data == "CreateTicketButton": 
			CreateTicketButton(upCQ.Message.Chat.ID, bot)			
		case upCQ.Data == "initiate":
			initiatebutton(update, bot)
	}
}

func CreateTicketName(chatID int64, bot *tgbotapi.BotAPI, language string){
	go help.ClearMessages1(chatID, bot)
	user, err := database.ReadUserByID(chatID)
	if user == nil {
		help.NewMessage(chatID, bot, fmt.Sprintf("%v",err), true)
		return
	}
	user.UserName = "............................................................................."
	user.Language = language
	user.CurrentTicket = 0
	user.MapUpdateOrCreate()
	help.NewMessage(chatID, bot, "Write how to address you", true)
	
	//go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	var defaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "пользователь")

	defaultKeyboard = [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("menu", "Menu"),
			tgbotapi.NewInlineKeyboardButtonData("createticket", "CreateTicketButton"),
		}, 
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

func initiatebutton(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	upCQ := update.CallbackQuery 
	db, err := database.OpenDB()
	if err != nil {
		fmt.Println(err)
		return 
	}
	if !database.IsTableExists(db, "staff"){
		upCQ := update.CallbackQuery.Message
		if err = database.CreateTable("staff"); err != nil {
			fmt.Println(err)
			help.NewMessage1(upCQ.Chat.ID, bot, fmt.Sprintf("%v", err), true)
		}
		staff := database.Staff{
			ChatID:				upCQ.Chat.ID,
			Admin:				1,	
			CurrentTicket: 		0,
			LinkName:			fmt.Sprintf("@%s",upCQ.Chat.UserName),
			UserName:			upCQ.Chat.FirstName,
			TicketClosed:		0,
			Rating:				0,
			Time: 		 		time.Now().Unix(),
		}
		if err := staff.InsertNew(); err != nil{
			help.NewMessage1(upCQ.Chat.ID, bot, fmt.Sprintf("Error initiating: %v", err), false)
		}
		staffbot.StartMenu(upCQ.Chat.ID, bot)
	} else {
		StartMenu(upCQ.Message.Chat.ID, bot)
	}
	db.Close()
}

func CreateTicketButton(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "Choose prefered language")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Russian", "LanguageRU"),
			tgbotapi.NewInlineKeyboardButtonData("English", "LanguageENG"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}
func CreateTicket(user database.User, bot *tgbotapi.BotAPI){
	ticketcr := database.Ticket{
		ChatID:			user.ChatID,
		SupChatID:		0,
		LinkName:		user.LinkName,
		SupLinkName:	"asd",
		UserName:		user.UserName,
		SupUserName:	"asd",
		Time: 			time.Now().Unix(),
		ClosingTime: 	0,
		Language:		user.Language,
		Status:			"Notificate",
	}
	if err := ticketcr.InsertNew(); err != nil{
		fmt.Println(err)
		help.NewMessage(user.ChatID, bot, fmt.Sprintf("%v", err), false)
		StartMenu(user.ChatID, bot)
		return
	}
	ticket, err := database.ReadOpenTicketByUserID(user.ChatID)
	if err != nil {
		fmt.Println(err)
		help.NewMessage(user.ChatID, bot, fmt.Sprintf("%v",err), false)
		return
	}
	user.CurrentTicket = ticket.TicketID
	user.MapUpdateOrCreate()
	if err = user.Update(); err != nil {
		fmt.Println(err)
		help.NewMessage(user.ChatID, bot, fmt.Sprintf("%v",err), false)
		return
	}

	go help.ClearMessages1(user.ChatID, bot)

	msg := tgbotapi.NewMessage(user.ChatID, "Ticket is open, You can ask.")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Close ticket", "TicketClose"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}