package staffbot

import (
	"fmt"
	"strconv"
	"time"
	"strings"

	database "tgbottrade/internal/database"
	help 	 "tgbottrade/pkg/api/help"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
func AdminStartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "admin panel")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("SupMenu", "Menu"),
			tgbotapi.NewInlineKeyboardButtonData("SupList", "SupList"),
			tgbotapi.NewInlineKeyboardButtonData("AddSup", "AddSup"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bookkeep", "BookkeepButton"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func AdminBookkeepMenu(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "Bookkeep menu. Find payment by:")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Invoice ID",  "BookkeepInvoiceButton"),
			tgbotapi.NewInlineKeyboardButtonData("Chat ID",  	"BookkeepChatIDButton"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", 	 	"adminMenu"),
			tgbotapi.NewInlineKeyboardButtonData("Date",	 	"BookkeepDateButton"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func AdminAddSupButton(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	staff.AddSup = true
	staff.MapUpdateOrCreate()
	msg := tgbotapi.NewMessage(chatID, "Write ChatID of the future support")
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
}

func AdminBackToMenuWithoutChanges(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	if value, exists := database.StaffMap[chatID]; exists{
		value.AddSup = false
		value.ChangeName = false
		value.FindByInvoice = false
		value.FindByChatID = false
		value.MapUpdateOrCreate()
	}
		AdminStartMenu(chatID, bot)
}

func BackToMenuWithoutChanges(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	staff.AddSup = false
	staff.ChangeName = false
	staff.MapUpdateOrCreate()
	Profile(chatID, bot, staff)
}

func AdminSupListButton(chatID int64, bot *tgbotapi.BotAPI){
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

func AdminSupProfile(chatID int64, bot *tgbotapi.BotAPI, supID string){
	go help.ClearMessages1(chatID, bot)
	id, _ := strconv.ParseInt(supID, 10, 64)
	staff, err := database.ReadStaffByID(id)
	if err != nil{
		fmt.Println(err)
		return
	}

	admin := "Support"
	if (staff.Admin > 0){
		admin = "Admin"
	}

	busy := "No"
	if (staff.CurrentTicket != 0){
		busy = strconv.FormatInt(staff.CurrentTicket, 10)
	}

	var defaultKeyboard []tgbotapi.InlineKeyboardButton
	defaultKeyboard = append(defaultKeyboard, tgbotapi.NewInlineKeyboardButtonData("Back", "adminMenu"))

	// Добавляем кнопку "Remove" только если chatID != idOutputPaymentsByID
	if chatID != id {
		defaultKeyboard = append(defaultKeyboard, tgbotapi.NewInlineKeyboardButtonData("Remove", fmt.Sprintf("RemoveButton%d", id)))
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n\nChatID: %d\nBusy: %s\nLink: %s\nName: %s\nClosed: %d\nRegistration: %s", admin, id, busy, staff.LinkName, staff.UserName, staff.TicketClosed, time.Unix(staff.Time, 0).Format("2006-01-02 15:04")))
	
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(defaultKeyboard)
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}

func AdminRemoveButton(chatID int64, bot *tgbotapi.BotAPI, supID string){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Are you sure you want to fire staff with ChatID:%s ?",supID))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", fmt.Sprintf("SupProfile%s",supID)),
				tgbotapi.NewInlineKeyboardButtonData("Remove", fmt.Sprintf("Remove%s",supID)),
			),
		)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}

func AdminRemove(chatID int64, bot *tgbotapi.BotAPI, supID string){
	go help.ClearMessages1(chatID, bot)

	id, _ := strconv.ParseInt(supID, 10, 64)

	if err := database.DeleteStaffByID(id); err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error removing staff with id: %s.\nError: %v", supID, err))
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Back to menu", "adminMenu"),
					tgbotapi.NewInlineKeyboardButtonData("Try again(In most cases it won't work)", fmt.Sprintf("Remove%s",supID)),
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
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Staff with id:%s successfully removed.", supID))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back to menu", "adminMenu"),
			),
		)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
	go help.ClearMessages1(id, bot)
}

func AdminBookkeepFindByInvoiceButton(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "Choose search option")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List InvoiceID", "BookkeepInvoiceList"),
			tgbotapi.NewInlineKeyboardButtonData("Write InvoiceID", "BookkeepInvoiceFind"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepButton"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)
}

func AdminBookkeepFindByInvoiceList(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	users, err := database.OutputInvoices()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepInvoiceButton"),
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

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "Payed invoice list")

	var row []tgbotapi.InlineKeyboardButton

	for _, user := range users{
		user := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d: %d", user.InvoiceID, user.ChatID), fmt.Sprintf("PaymentInvoiceID%d", user.InvoiceID))
		row = append(row, user)

		if len(row) == 2 {
			DefaultKeyboard = append(DefaultKeyboard, row)	
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		DefaultKeyboard = append(DefaultKeyboard, row)
	}

	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepInvoiceButton"),
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

func AdminBookkeepFindByInvoiceFindButton(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	staff.FindByInvoice = true
	staff.MapUpdateOrCreate()

	msg := tgbotapi.NewMessage(chatID, "Write the invoice ID which you want to check")
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
}
func AdminBookkeepFindByChatIDFind(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	go help.ClearMessages1(chatID, bot)
	staff.FindByChatID = true
	staff.MapUpdateOrCreate()

	msg := tgbotapi.NewMessage(chatID, "Write the chat ID which you want to check")
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
}
func AdminPaymentInfoID(chatID int64, bot *tgbotapi.BotAPI, invoiceID string){
	//go help.ClearMessages1(chatID, bot)
	id, _ := strconv.ParseInt(invoiceID, 10, 64)

	payment, err := database.OutputPaymentByInvoiceID(id)
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error: %v", err), true)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Payment info\n\nInvoiceID: %d\nChatID: %d\nLink: %s\nAmount: %d\nAsset: %s\nTime: %s", payment.InvoiceID, payment.ChatID, payment.LinkName, payment.Amount, payment.Asset, time.Unix(payment.PaymentTime, 0).Format("2006-01-02 15:04")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", fmt.Sprintf("PaymentChatID%d",payment.ChatID)),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func AdminBookkeepFindByChatIDButton(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	msg := tgbotapi.NewMessage(chatID, "Choose search option")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List ChatID", "BookkeepChatIDList"),
			tgbotapi.NewInlineKeyboardButtonData("Write ChatID", "BookkeepChatIDFind"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepButton"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}
func AdminBookkeepFindByChatIDList(chatID int64, bot *tgbotapi.BotAPI){
	go help.ClearMessages1(chatID, bot)

	users, err := database.OutputPayedIDs()
	fmt.Println(users)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "adminMenu"),
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

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "Payed ID's list")

	var row []tgbotapi.InlineKeyboardButton

	for _, user := range users{
		user := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d",user.ChatID), fmt.Sprintf("PaymentChatID%d", user.ChatID))
		row = append(row, user)

		if len(row) == 2 {
			DefaultKeyboard = append(DefaultKeyboard, row)	
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		DefaultKeyboard = append(DefaultKeyboard, row)
	}

	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepChatIDButton"),
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

func AdminPaymentInfoUserListID(chatID int64, bot *tgbotapi.BotAPI, idd string){
	help.ClearMessages1(chatID, bot)
	id, _ := strconv.ParseInt(idd, 10, 64)
	payments, err := database.OutputPaymentsByID(id)
	if err != nil { 
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "AdminBackToMenuWithoutChanges"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, er := bot.Send(msg)
		if er != nil {
			fmt.Println("Error sending start menu: ", er)
		}
		go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
		return
	}
	if value, exists := database.StaffMap[chatID]; exists{
		value.FindByChatID = false
		value.MapUpdateOrCreate()
	}
	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("User payment list with ID: %d", id))

	for _, payment := range payments{
		payment1 := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d %s : %s",payment.Amount, payment.Asset, time.Unix(payment.PaymentTime, 0).Format("2006.01.02 15:04")) , fmt.Sprintf("PaymentID%d", payment.InvoiceID)),
		}
		DefaultKeyboard = append(DefaultKeyboard, payment1)					
	}
			
	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepChatIDList"),
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

func AdminBookkeepFindByInvoiceID(chatID int64, bot *tgbotapi.BotAPI, idd string){
	id, _ := strconv.ParseInt(idd, 10, 64)

	payment, err := database.OutputPaymentByInvoiceID(id)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepInvoiceButton"),
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

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Payment info\n\nInvoiceID: %d\nChatID: %d\nLink: %s\nAmount: %d\nAsset: %s\nTime: %s", payment.InvoiceID, payment.ChatID, payment.LinkName, payment.Amount, payment.Asset, time.Unix(payment.PaymentTime, 0).Format("2006-01-02 15:04")))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepInvoiceButton"),
		),
	)
	msg.ReplyMarkup = keyboard
	sent, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Error sending start menu: ", err)
	}
	go help.AddToDelete1(sent.Chat.ID, sent.MessageID)	
}

func AdminBookkeepDateButton(chatID int64, bot *tgbotapi.BotAPI){
	help.ClearMessages1(chatID, bot)
	invoices, err := database.OutputInvoices()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error outputing invoices: %v", err), true)
		return
	}
	var years []int
	var lyear	int
	count := 0 
	for _, invoice := range invoices{
		count++
		year, _, _ := time.Unix(invoice.PaymentTime, 0).Date()
		if count == 1 {
			lyear = year
			years = append(years, year)
			continue
		}
		if lyear == year{
			continue
		} else {
			lyear = year
			years = append(years, year)
		}
	}

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, "Years with profit:")

	for _, newyear := range years{
		year1 := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", newyear), fmt.Sprintf("DateListMonthOf%d", newyear)),
		}
		DefaultKeyboard = append(DefaultKeyboard, year1)					
	}
			
	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepButton"),
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

func AdminBookkeepDateListMonth(chatID int64, bot *tgbotapi.BotAPI, yea string){
	year, _ := strconv.Atoi(yea)
	if year == 0 {
		fmt.Println("Error convertation")
		return
	}
	invoices, err := database.OutputInvoices()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error outputing invoices: %v", err), true)
		return
	}

	var months []int
	var lmonth	int
	count := 0 
	for _, invoice := range invoices{
		yearr, month, _ := time.Unix(invoice.PaymentTime, 0).Date()
		if year != yearr{
			continue
		}
		count++
		if count == 1 {
			lmonth = int(month)
			months = append(months, int(month))
			continue
		}
		if lmonth == int(month){
			continue
		} else {
			lmonth = int(month)
			months = append(months, int(month))
		}
	}

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Months with profit in %s:", yea))

	for _, newmonth := range months{
		month1 := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)",time.Month(newmonth), newmonth), fmt.Sprintf("DateListDayOf%s:%d", yea, newmonth)),
		}
		DefaultKeyboard = append(DefaultKeyboard, month1)					
	}
			
	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", "BookkeepDateButton"),
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

func AdminBookkeepDateListDay(chatID int64, bot *tgbotapi.BotAPI, data string){
	confirmeddata := strings.Split(data, ":")
	year, _ := strconv.Atoi(confirmeddata[0]) 
	month, _ := strconv.Atoi(confirmeddata[1])

	invoices, err := database.OutputInvoices()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error outputing invoices: %v", err), true)
		return
	}

	var days []int
	var lday   int
	count := 0 
	for _, invoice := range invoices{
		year1, month1, day := time.Unix(invoice.PaymentTime, 0).Date()
		if year1 != year && month != int(month1){
			continue
		}
		count++
		if count == 1 {
			lday = day
			days = append(days, day)
			continue
		}
		if lday == day{
			continue
		} else {
			lday = day
			days = append(days, day)
		}
	}

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Days with profit in %d.%d:", year, month))

	for _, newday := range days{
		day1 := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", newday), fmt.Sprintf("DateListInvoiceOf%d:%d:%d", year, month, newday)),
		}
		DefaultKeyboard = append(DefaultKeyboard, day1)					
	}
			
	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", fmt.Sprintf("DateListMonthOf%d",year)),
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

func AdminBookkeepDateListInvoice(chatID int64, bot *tgbotapi.BotAPI, data string){
	confirmeddata := strings.Split(data, ":")
	year, _ := strconv.Atoi(confirmeddata[0]) 
	month, _ := strconv.Atoi(confirmeddata[1])
	day, _ := strconv.Atoi(confirmeddata[2])

	invoices, err := database.OutputInvoices()
	if err != nil {
		help.NewMessage1(chatID, bot, fmt.Sprintf("Error outputing invoices: %v", err), true)
		return
	}

	var DefaultKeyboard [][]tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Profit in %d.%d.%d:", year, month, day))

	for _, invoice := range invoices{
		year1, month1, day1 := time.Unix(invoice.PaymentTime, 0).Date()
		if year1 != year && month != int(month1) && day != day1{
			continue
		}
		hour, min, _ := time.Unix(invoice.PaymentTime, 0).Clock()
		pay := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d:%d", hour, min), fmt.Sprintf("PaymentID%d", invoice.InvoiceID)),
		}
		DefaultKeyboard = append(DefaultKeyboard, pay)	
	}

	back := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Back", fmt.Sprintf("DateListDayOf%s",data)),
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