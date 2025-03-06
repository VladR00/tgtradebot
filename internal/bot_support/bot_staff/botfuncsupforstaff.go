package staffbot

import (
	"fmt"
	"strings"
	"strconv"

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
	switch {
		case strings.HasPrefix(upCQ.Data, "Accept"):
			AcceptTicket(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Accept"))
		case upCQ.Data == "Menu":
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
func NotificateSups(user database.User, bot *tgbotapi.BotAPI){
	stafflist, err := database.OutputStaffWithCurrTicketNull()
	if err != nil {
		fmt.Println(err)
		return
	} 
	for _, staff := range stafflist{
		msg := tgbotapi.NewMessage(staff.ChatID, fmt.Sprintf("New ticket with ID: %d\nUsername: %s\nPrefered language: %s\n", user.CurrentTicket, user.UserName, user.Language))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Accept", fmt.Sprintf("Accept%d",user.CurrentTicket)),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
	}
	ticket, err := database.ReadTicketByID(user.CurrentTicket)
	if err != nil {
		fmt.Println(err)
	}
	ticket.Status = "Open"
	fmt.Println(ticket)
	err = ticket.Update()
	if err != nil {
		fmt.Println(err)
	}
}

func AcceptTicket(chatID int64, bot *tgbotapi.BotAPI, ticketid string){
	t, _ := strconv.ParseInt(ticketid, 10, 64)
	ticketcr, err := database.ReadTicketByID(t)
	if err != nil{
		fmt.Println(err)
		return
	}
	if ticketcr.SupChatID != 0 {
		help.NewMessage(chatID, bot, fmt.Sprintf("Ticket %d has already been taken by: %s", ticketcr.TicketID, ticketcr.UserName), true)
		return
	}
	staff, err := database.ReadStaffByID(chatID)
	if err != nil{
		fmt.Println(err)
		return
	}
	//editmessageforall()
	ticket := database.Ticket{
		TicketID:		ticketcr.TicketID,
		ChatID:			ticketcr.ChatID,
		SupChatID:		staff.ChatID,
		LinkName:		ticketcr.LinkName,
		SupLinkName:	staff.LinkName,
		UserName:		ticketcr.UserName,
		SupUserName:	staff.UserName,
		Time: 			ticketcr.Time,
		ClosingTime: 	0,
		Language:		ticketcr.Language,
		Status:			"Chating",
	}
	ticket.Update()
	staff.CurrentTicket = ticketcr.TicketID
	staff.MapUpdateOrCreate()
	staff.Update()
	messages, err := database.ReadAllMessagesByTicketID(ticket.TicketID)
	if err != nil {
		fmt.Println(err)
		return
	}
	SendAllMessages(messages, ticket, bot)
}
func SendAllMessages(messages []*database.TicketMessage, ticket database.Ticket, bot *tgbotapi.BotAPI){
	for _, message := range messages{
		msg := tgbotapi.NewCopyMessage(ticket.SupChatID, message.ChatID, message.MessageID)
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending: ", err)
		}
		go help.AddToDelete1(ticket.SupChatID, sent.MessageID)
	}
}