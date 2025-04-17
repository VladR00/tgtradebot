package staffbot

import (
	"fmt"
	"strconv"
	"time"

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
			tgbotapi.NewInlineKeyboardButtonData("Invoice ID", "BookkeepInvoiceFind"),
			tgbotapi.NewInlineKeyboardButtonData("Date",	   "BookkeepDateButton"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", 	 "adminMenu"),
			tgbotapi.NewInlineKeyboardButtonData("Chat ID",  "BookkeepChatIDButton"),
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

func AdminBackToMenuWithoutChanges(chatID int64, bot *tgbotapi.BotAPI, staff *database.Staff){
	staff.AddSup = false
	staff.ChangeName = false
	staff.MapUpdateOrCreate()
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

	// Добавляем кнопку "Remove" только если chatID != id
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