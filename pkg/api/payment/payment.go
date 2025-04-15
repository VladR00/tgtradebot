package payment

import (
	"time"
	"fmt"
	"log"
	"strconv"
	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"

	help	 "tgbottrade/pkg/api/help"
	database "tgbottrade/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TopUp(bot *tgbotapi.BotAPI, chatID int64, client *cryptobot.Client, asset string, amount string) {
	go help.ClearMessages(chatID, bot)
	help.NewMessage(chatID, bot, "Creating invoice...",  true)
	invoice, err := client.CreateInvoice(cryptobot.CreateInvoiceRequest{
		Asset:          asset,
		Amount:         amount,
		Description:    "Successfully Paid",
		HiddenMessage:  "Thank you for using us",
		PaidBtnName:    "viewItem",
		PaidBtnUrl:     "https://t.me/CryptoTestnetBot?start=IV6UgcC5mlW3",
		Payload:        "any payload we need in our application",
		AllowComments:  false,
		AllowAnonymous: false,
		ExpiresIn:      60 * 5,
	})
	if err != nil {
		log.Println("Error creating invoice:", err)
		help.NewMessage(chatID, bot, "Error creating invoice",  true)
		go func() {
			time.Sleep(3 * time.Second)
			TopUp(bot, chatID, client, asset, amount)
		}()
	}
	go help.ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Pay %s %s to address %s", invoice.Amount, invoice.Asset, invoice.PayUrl))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Close payment", "Menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
	go help.AddToDelete(sentMsg.Chat.ID, sentMsg.MessageID)
	go CheckPaymentStatus(bot, chatID, client, strconv.FormatInt(invoice.ID, 10)+strconv.FormatInt(chatID, 10))
}
func CheckPaymentStatus(bot *tgbotapi.BotAPI, chatID int64, client *cryptobot.Client, targetinvoice string) {
	for i := 0; i < 30; i++ {
		invoices, err := client.GetInvoices(nil)
		if err != nil {
			log.Println("Error getting invoices:", err)
			return
		}

		for _, invoice := range invoices {
			if invoice.Status == cryptobot.InvoicePaidStatus {
				if strconv.FormatInt(invoice.ID, 10)+strconv.FormatInt(chatID, 10) == targetinvoice {
					go help.ClearMessages(chatID, bot)

					topup, err := strconv.Atoi(invoice.Amount)
					if err != nil {
						fmt.Printf("Error convert: %w\n", err)
					}

					user, _ := database.ReadUserByID(chatID)
					invoicestruct := database.Invoice{
						InvoiceID:		invoice.ID,
						ChatID:			chatID,
						LinkName:		user.LinkName,
						Amount:			int64(topup),
						StringAmount:	invoice.Amount,
						Asset:			invoice.Asset,
						PaymentTime:	time.Now().Unix(),
					}
					
					if err := invoicestruct.InsertNew(); err != nil {
						help.NewMessage(chatID, bot, fmt.Sprintf("Error insert new payment with ID:%d for ChatID:%d with LinkName:%s. Amount:%d (string:%s), Asset:%s, Time:%s\n\nContact us with forward this message.\n\nError: %v", invoicestruct.InvoiceID, chatID, invoicestruct.LinkName, invoicestruct.Amount, invoicestruct.StringAmount, time.Unix(invoicestruct.PaymentTime, 0).Format("2006-01-02 15:04"), err), false)
					}

					if user != nil{
						user.Balance = user.Balance + int64(topup)
						if err := user.Update(); err != nil{
							help.NewMessage(chatID, bot, fmt.Sprintf("Error update balance (+%d). Contact us with forward this message.\nError: %v", int64(topup), err), false)
						}
					}
					msg := tgbotapi.NewMessage(chatID, "good")
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Menu", "Menu"),
						),
					)
					msg.ReplyMarkup = keyboard
					sentMsg, err := bot.Send(msg)
					if err != nil {
						log.Println("Error sending message: ", err)
					}
					go help.AddToDelete(sentMsg.Chat.ID, sentMsg.MessageID)
					log.Printf("Invoice %d paid!\n", invoice.ID)
					help.NewMessage(chatID, bot, fmt.Sprintf("Invoice %d paid!\n", invoice.ID), false)
					return
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}