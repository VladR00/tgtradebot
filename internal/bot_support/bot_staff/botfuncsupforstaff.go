package staffbot

import (
	"fmt"
	"strings"
	"strconv"
	"time"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessageSwitchForAuthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upM := update.Message
	fmt.Printf("Handle message on support bot from staff: %s. From user: %s\n", upM.Text, upM.Chat.UserName)
	if staff.LinkName == "Agent"{
		staff.LinkName = fmt.Sprintf("@%s",upM.Chat.UserName)
		staff.UserName = upM.Chat.FirstName
		staff.Update()
		if _, exists := database.StaffMap[upM.Chat.ID]; exists{ 
			staff.MapUpdateOrCreate()
		}
	}
	fmt.Println(database.UserMap)
	fmt.Println(database.StaffMap)
	fmt.Println(database.TicketMap)
	if value, exists := database.StaffMap[upM.Chat.ID]; exists{
		if (value.CurrentTicket != 0){
			message := database.TicketMessage{
				TicketID:	value.CurrentTicket,
				Support:	1,
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
			msg := tgbotapi.NewCopyMessage(ticket.ChatID, message.ChatID, message.MessageID)
			sent, err := bot.Send(msg)
			if err != nil {
				fmt.Println("Error sending: ", err)
			} else {
				go help.AddToDelete1(ticket.ChatID, sent.MessageID)
			}
			return
		} else if (value.AddSup){
			id, _ := strconv.ParseInt(upM.Text, 10, 64)
			if id == 0 {
				help.NewMessage1(upM.Chat.ID, bot, fmt.Sprintf("Wrong ChatID: %d. Try another one", id), true)
				return
			}
			newstaff := database.Staff{
				ChatID:				id,
				Admin:				0,
				CurrentTicket: 		0,
				LinkName:			"Agent",
				UserName:			"nil",
				TicketClosed:		0,
				Rating:				0,
				Time: 		 		time.Now().Unix(),
			}
			if err := newstaff.InsertNew(); err != nil{
				help.NewMessage1(upM.Chat.ID, bot, fmt.Sprintf("Error initiating: %v", err), false)
				return
			}
			msg := tgbotapi.NewMessage(upM.Chat.ID, fmt.Sprintf("Support successfully added by ID: %d", id))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Back", "AdminBackToMenuWithoutChanges"),
				),
			)
			msg.ReplyMarkup = keyboard
			sent, err := bot.Send(msg)
			if err != nil {
				fmt.Println("Error sending start menu: ", err)
			}
			go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
			return
		} else if (value.ChangeName){
			value.UserName = upM.Text
			value.MapUpdateOrCreate()
			value.Update()
			msg := tgbotapi.NewMessage(upM.Chat.ID, fmt.Sprintf("New name confirmed: %s", value.UserName))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Back", "BackToProfile"),
				),
			)
			msg.ReplyMarkup = keyboard
			sent, err := bot.Send(msg)
			if err != nil {
				fmt.Println("Error sending start menu: ", err)
			}
			go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
			return
		} else if (value.FindByInvoice){
			go help.ClearMessages1(upM.Chat.ID, bot)
			id, _ := strconv.ParseInt(upM.Text, 10, 64)
			
			if id == 0 {
				msg := tgbotapi.NewMessage(upM.Chat.ID, fmt.Sprintf("Wrong InvoiceID: %d. Try another one", id))
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Back", "AdminBackToMenuWithoutChanges"),
					),
				)
				msg.ReplyMarkup = keyboard
				sent, err := bot.Send(msg)
				if err != nil {
					fmt.Println("Error sending start menu: ", err)
				}
				go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
				return
			}
			AdminBookkeepFindByInvoiceID(upM.Chat.ID, bot, upM.Text)
			return
		} else if (value.FindByChatID) {
			go help.ClearMessages1(upM.Chat.ID, bot)
			id, _ := strconv.ParseInt(upM.Text, 10, 64)
			
			if id == 0 {
				msg := tgbotapi.NewMessage(upM.Chat.ID, fmt.Sprintf("Wrong ChatID: %d. Try another one", id))
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Back", "AdminBackToMenuWithoutChanges"),
					),
				)
				msg.ReplyMarkup = keyboard
				sent, err := bot.Send(msg)
				if err != nil {
					fmt.Println("Error sending start menu: ", err)
				}
				go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
				return
			}

			AdminPaymentInfoUserListID(upM.Chat.ID, bot, upM.Text)
		}
	}
	
	switch upM.Text {
		case "/start":
			StartMenu(upM.Chat.ID, bot, staff)
	}
}

func HandleCallBackSwitchForAuthorizedInTableStaff(update tgbotapi.Update, bot *tgbotapi.BotAPI, staff *database.Staff){
	upCQ := update.CallbackQuery
	fmt.Printf("Handle callback on support bot from staff: %s. From user: %s\n", upCQ.Data, upCQ.Message.Chat.UserName)
	if staff.LinkName == "Agent"{
		staff.LinkName = fmt.Sprintf("@%s",upCQ.Message.Chat.UserName)
		staff.UserName = upCQ.Message.Chat.FirstName
		staff.Update()
		if _, exists := database.StaffMap[upCQ.Message.Chat.ID]; exists{ 
			staff.MapUpdateOrCreate()
		}
	}
	switch upCQ.Data {
		case "Menu":
			StartMenu(upCQ.Message.Chat.ID, bot, staff)
			return
		case "adminMenu":
			AdminStartMenu(upCQ.Message.Chat.ID, bot)
			return
		case "AddSup":
			AdminAddSupButton(upCQ.Message.Chat.ID, bot, staff)
			return 
		case "AdminBackToMenuWithoutChanges": 
			AdminBackToMenuWithoutChanges(upCQ.Message.Chat.ID, bot, staff)
			return
		case "BackToProfile": 
			BackToMenuWithoutChanges(upCQ.Message.Chat.ID, bot, staff)
			return
		case "SupList":
			AdminSupListButton(upCQ.Message.Chat.ID, bot)
			return 
		case "Profile":
			Profile(upCQ.Message.Chat.ID, bot, staff)
			return
		case "Change name":
			ChangeName(upCQ.Message.Chat.ID, bot, staff)
			return 
		case "BookkeepButton":
			AdminBookkeepMenu(upCQ.Message.Chat.ID, bot)
			return
		case "BookkeepInvoiceButton":
			AdminBookkeepFindByInvoiceButton(upCQ.Message.Chat.ID, bot)
			return
		case "BookkeepInvoiceFind":
			AdminBookkeepFindByInvoiceFindButton(upCQ.Message.Chat.ID, bot, staff)
			return
		case "BookkeepInvoiceList":
			AdminBookkeepFindByInvoiceList(upCQ.Message.Chat.ID, bot)
			return
		case "BookkeepChatIDButton":
			AdminBookkeepFindByChatIDButton(upCQ.Message.Chat.ID, bot)
			return
		case "BookkeepChatIDFind":
			AdminBookkeepFindByChatIDFind(upCQ.Message.Chat.ID, bot, staff)
			return
		case "BookkeepChatIDList":
			AdminBookkeepFindByChatIDList(upCQ.Message.Chat.ID, bot)
			return 
		case "BookkeepDateButton":
			AdminBookkeepDateButton(upCQ.Message.Chat.ID, bot)
			return
		case "Turn aside":
			TurnAside(upCQ.Message.Chat.ID, bot, staff)
	}

	switch {
		case strings.HasPrefix(upCQ.Data, "Accept"):
			AcceptTicket(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Accept"), staff)
		case strings.HasPrefix(upCQ.Data, "Close"):
			CloseTicket(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Close"))
		case strings.HasPrefix(upCQ.Data, "SupProfile"):
			AdminSupProfile(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "SupProfile"))
		case strings.HasPrefix(upCQ.Data, "RemoveButton"):
			AdminRemoveButton(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "RemoveButton"))
		case strings.HasPrefix(upCQ.Data, "Remove"):
			AdminRemove(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "Remove"))
		case strings.HasPrefix(upCQ.Data, "PaymentID"): 
			AdminPaymentInfoID(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "PaymentID"))
		case strings.HasPrefix(upCQ.Data, "PaymentChatID"): 
			AdminPaymentInfoUserListID(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "PaymentChatID"))
		case strings.HasPrefix(upCQ.Data, "PaymentInvoiceID"): 
			AdminBookkeepFindByInvoiceID(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "PaymentInvoiceID"))
		case strings.HasPrefix(upCQ.Data, "DateListMonthOf"):
			AdminBookkeepDateListMonth(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "DateListMonthOf"))
		case strings.HasPrefix(upCQ.Data, "DateListDayOf"):
			AdminBookkeepDateListDay(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "DateListDayOf"))
		case strings.HasPrefix(upCQ.Data, "DateListInvoiceOf"):
			AdminBookkeepDateListInvoice(upCQ.Message.Chat.ID, bot, strings.TrimPrefix(upCQ.Data, "DateListInvoiceOf"))
		case strings.HasPrefix(upCQ.Data, "Move"):
			MoveToClaimedTicket(upCQ.Message.Chat.ID, bot, staff, strings.TrimPrefix(upCQ.Data, "Move"))
	}
}	

func MoveToClaimedTicket(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff, ticccket string){
	id, _ := strconv.ParseInt(ticccket, 10, 64)
	if id == 0 {
		fmt.Println("Error convert ticket id")
		return
	}

	ticket, err := database.ReadTicketByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	ViewTicket(chatID, bot, *ticket, staff)
}

func StartMenu(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	var defaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "support")

	defaultKeyboard = [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Menu", "Menu"),
			tgbotapi.NewInlineKeyboardButtonData("Profile", "Profile"),
		}, 
	}

	adminpanel := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("admin panel", "adminMenu"),
	}

	if staff.Admin == 1{
		defaultKeyboard = append(defaultKeyboard, adminpanel)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(defaultKeyboard...)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
	go ViewOpenTickets(chatID, bot, staff)
}

func Profile(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ChatID: %d\nCurrent ticket: %d\nLink: %s\nName: %s\nClosed tickets: %d\nRegistration: %s", staff.ChatID, staff.CurrentTicket, staff.LinkName, staff.UserName, staff.TicketClosed, time.Unix(staff.Time, 0).Format("2006-01-02 15:04")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "Menu"),
			tgbotapi.NewInlineKeyboardButtonData("Change name", "Change name"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}

func ChangeName(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	if value, exists := database.StaffMap[chatID]; exists{
		value.ChangeName = true
		value.MapUpdateOrCreate()
		go help.ClearMessages1(chatID, bot)
		help.NewMessage1(chatID, bot, "You can write your new name", true)
	} else {
		staff.ChangeName = true
		staff.MapUpdateOrCreate()
		go help.ClearMessages1(chatID, bot)
		help.NewMessage1(chatID, bot, "You can write your new name", true)
	}
	return
}

func AcceptTicket(chatID int64, bot *tgbotapi.BotAPI, ticketid string, staff *database.Staff){
	t, _ := strconv.ParseInt(ticketid, 10, 64)
	ticketcr, err := database.ReadTicketByID(t)
	if err != nil{
		fmt.Println(err)
		return
	}
	if ticketcr.SupChatID != 0 && ticketcr.SupChatID != chatID{
		help.NewMessage1(chatID, bot, fmt.Sprintf("Ticket %d has already been taken by: %s", ticketcr.TicketID, ticketcr.SupUserName), true)
		go func(){
			time.Sleep(1 * time.Second)
			StartMenu(chatID, bot, staff)
		}()
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
		Status:			"Chat",
	}
	ticket.Update()
	ticket.MapUpdateOrCreate()
	ViewTicket(chatID, bot, ticket, staff)
}
func ViewTicket(chatID int64, bot *tgbotapi.BotAPI, ticket database.Ticket, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	staff.CurrentTicket = ticket.TicketID
	staff.MapUpdateOrCreate()
	staff.Update()
	payments, err := database.OutputPaymentsByID(ticket.ChatID)
	if err != nil {
		fmt.Println(err)
	}	
	var resultout string
	countpayments := 0
	var asset string
	var amount string
	asdfg := false 

	if payments == nil {
		resultout = "\nPayments: 0"
	} else {
		for _, payment := range payments{
			countpayments++
			if countpayments == 1 {
				asset = payment.Asset
				amount = payment.StringAmount
				continue
			}
			if asset != payment.Asset{
				asdfg = true
				asset = fmt.Sprintf("%s, %s", asset, payment.Asset)
				amount = fmt.Sprintf("%s + %s", amount, payment.StringAmount)
				continue
			}		
			if !asdfg{	
				qwer, _ := strconv.ParseInt(amount, 10, 64) 
				amount = strconv.FormatInt(qwer + payment.Amount, 10)
			} else {
				asset = fmt.Sprintf("%s, %s", asset, payment.Asset)
				amount = fmt.Sprintf("%s + %s", amount, payment.StringAmount)
			}
		}
		resultout = fmt.Sprintf("\nPayments:%d\nSpent:%s\nAssets:%s", countpayments, amount, asset)
	}

	msg := tgbotapi.NewMessage(staff.ChatID, fmt.Sprintf("Ticket info\nID: %d\nName: %s\nPrefered language: %s%s\nOpen time: %s", ticket.TicketID, ticket.UserName, ticket.Language, resultout, time.Unix(ticket.Time, 0).Format("2006-01-02 15:04")))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Close", fmt.Sprintf("Close%d",ticket.TicketID)),
				tgbotapi.NewInlineKeyboardButtonData("Turn aside", "Turn aside"),
			),
		)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
	SendAllMessages(ticket, bot)
}

func SendAllMessages(ticket database.Ticket, bot *tgbotapi.BotAPI){
	messages, err := database.ReadAllMessagesByTicketID(ticket.TicketID)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, message := range messages{
		msg := tgbotapi.NewCopyMessage(ticket.SupChatID, message.ChatID, message.MessageID)
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending: ", err)
		}
		go help.AddToDelete1(ticket.SupChatID, sent.MessageID)
	}
}
func CloseTicket(chatID int64, bot *tgbotapi.BotAPI, ticketid string){
	t, _ := strconv.ParseInt(ticketid, 10, 64)
	ticketcr, err := database.ReadTicketByID(t)
	if err != nil{
		fmt.Println(err)
		return
	}
	staff, err := database.ReadStaffByID(chatID)
	if err != nil{
		fmt.Println(err)
		return
	}
	ticket := database.Ticket{
		TicketID:		ticketcr.TicketID,
		ChatID:			ticketcr.ChatID,
		SupChatID:		ticketcr.SupChatID,
		LinkName:		ticketcr.LinkName,
		SupLinkName:	ticketcr.SupLinkName,
		UserName:		ticketcr.UserName,
		SupUserName:	ticketcr.SupUserName,
		Time: 			ticketcr.Time,
		ClosingTime: 	time.Now().Unix(),
		Language:		ticketcr.Language,
		Status:			"Closed",
	}
	staff.TicketClosed++
	staff.CurrentTicket = 0
	ticket.Update()
	ticket.MapDelete()
	staff.MapDelete()
	staff.Update()
	if value, exists := database.UserMap[ticket.ChatID]; exists{
		value.MapDelete()
		msg := tgbotapi.NewMessage(ticket.ChatID, fmt.Sprintf("Your ticket with ID: %d is closed", ticket.TicketID))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Menu", "Menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		} else {
			go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
		}
	}
	go help.ClearMessages1(chatID, bot)
	StartMenu(chatID, bot, staff)
}

func ViewOpenTickets(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	tickets, err := database.ReadOpenTickets()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Tickets can't load:%v", err), true)
	}

	var ClaimedKeyboard [][]tgbotapi.InlineKeyboardButton

	//var row []tgbotapi.InlineKeyboardButton
	for _, ticket := range tickets{
		if (ticket.Status == "Open"){
			continue
		}

		if (ticket.Status == "Chat" || ticket.Status == "Claimed" && ticket.SupChatID == chatID){
			user := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d",ticket.TicketID), fmt.Sprintf("Move%d",ticket.TicketID))
			ClaimedKeyboard = append(ClaimedKeyboard, tgbotapi.NewInlineKeyboardRow(user))	
		}
	}
	if len(ClaimedKeyboard) > 0 {
		msg := tgbotapi.NewMessage(chatID, "Claimed tickets by ID:")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(ClaimedKeyboard...)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: приколы тут:", err)
		}
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
	}
	for _, ticket := range tickets{
		if ticket.Status == "Chat" || ticket.Status == "Claimed"{
			continue
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ticket ID:%d, Language:%s\nUserName:%s\nOpen time:%s", ticket.TicketID, ticket.Language, ticket.UserName, time.Unix(ticket.Time, 0).Format("2006-01-02 15:04")))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Accept", fmt.Sprintf("Accept%d",ticket.TicketID)),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending start menu: ", err)
		}
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
	}
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

func TurnAside(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	// if value, exists := database.StaffMap[chatID]; exists{
	// 	staff = value 
	// }
	go help.ClearMessages1(chatID, bot)
	var ticket database.Ticket
	if value, exists := database.TicketMap[staff.CurrentTicket]; exists{
		ticket = value 
	}
	staff.CurrentTicket = 0
	staff.Update()
	staff.MapUpdateOrCreate()

	ticket.Status = "Claimed"
	ticket.Update()
	ticket.MapUpdateOrCreate()

	StartMenu(chatID, bot, staff)
}