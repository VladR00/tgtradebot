package main

import (
	"log"
	"os"
	"fmt"
	"sync"
	"bufio"
	"strconv"
	"strings"
	"time"
	"database/sql"
	_ "modernc.org/sqlite"

	"github.com/arthurshafikov/cryptobot-sdk-golang/cryptobot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
var (
	messagesMutex   sync.Mutex
	dbMutex			sync.Mutex
	dbpath 			string 		=	"./db.sqlite"
)

// func ExistingTable(tableName string) bool{
// 	dbMutex.Lock()
// 	defer dbMutex.Unlock()

// 	_, err = os.Stat(dbpath)
// 	if os.IsNotExist(err) {
// 		log.Println("first creating file db.")
// 	}

// 	db, err := sql.Open("sqlite", dbpath)
// 	if err != nil {
// 		log.Println("Error open/create DB: ",err)
// 	}
// 	defer db.Close()

// 	var exists int
// 	query := fmt.Sprintf(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%s';`, tableName)
// 	err = db.QueryRow(query).Scan(&exists)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Println("No rows returned, table does not exist.");
// 			return false;
// 		} else {
// 			log.Println("Error sql request: ", err);
// 			return false;
// 		}
// 		return false;
// 	}
// 	if exist > 0 {
// 		log.Println("Table exist");
// 		return true;
// 	}
// 	log.Println("Understandble error")
// 	return false;
// }

func DBCheckExisting() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		log.Println("Error open/create DB: ",err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (chatid INTEGER UNIQUE, name TEXT, language TEXT, balance INTEGER, date DATETIME)`)
	if err != nil {
		log.Println("Error sql request create users: ",err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS support (chatid INTEGER UNIQUE, name TEXT, date DATETIME)`)
	if err != nil {
		log.Println("Error sql request create support: ",err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS token (tokenid INTEGER PRIMARY KEY, userchatId INTEGER, supportchatID INTEGER, language TEXT, dateCreate DATETIME, dateCreateAccount, dateBuy DATETIME)`)
	if err != nil {
		log.Println("Error sql request create tokens: ",err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokenmessages (tokenid INTEGER, message STRING, date DATETIME)`)
	if err != nil {
		log.Println("Error sql request create tokenmessages: ",err)
	}
}

func StartMenu(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "hello muchahos")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Profile", "Profile"),
				tgbotapi.NewInlineKeyboardButtonData("Services", "Services"),
				tgbotapi.NewInlineKeyboardButtonData("F.A.Q.", "FAQ"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}

func NewMessage(chatID int64, bot *tgbotapi.BotAPI,message string, needDelete bool){
	sent, err := bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		log.Println("Error sending message: ", err)
	}
	if (needDelete){
		go AddToDelete(sent.Chat.ID, sent.MessageID)
	}
}
func ServiceMenu(chatID int64, bot *tgbotapi.BotAPI){
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, "Services")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("anal 1000", "topup1000"),
				tgbotapi.NewInlineKeyboardButtonData("blow 500", "topup500"),
				tgbotapi.NewInlineKeyboardButtonData("dance 200", "topup200"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Back", "startmenu"),
			),
		)
		msg.ReplyMarkup = keyboard
		sent, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending start menu: ", err)
		}
		go AddToDelete(sent.Chat.ID, sent.MessageID)	
}

func TopUp(bot *tgbotapi.BotAPI, chatID int64, client *cryptobot.Client, asset string, amount string) {
	go ClearMessages(chatID, bot)
	NewMessage(chatID, bot, "Creating invoice...",  true)
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
		NewMessage(chatID, bot, "Error creating invoice",  true)
		go func() {
			time.Sleep(3 * time.Second)
			TopUp(bot, chatID, client, asset, amount)
		}()
	}
	go ClearMessages(chatID, bot)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Pay %s %s to address %s", invoice.Amount, invoice.Asset, invoice.PayUrl))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Close payment", "closeInvoice"),
		),
	)
	msg.ReplyMarkup = keyboard
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
	go AddToDelete(sentMsg.Chat.ID, sentMsg.MessageID)
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
					go ClearMessages(chatID, bot)
					msg := tgbotapi.NewMessage(chatID, "good")
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Menu", "startmenu"),
						),
					)
					msg.ReplyMarkup = keyboard
					sentMsg, err := bot.Send(msg)
					if err != nil {
						log.Println("Error sending message: ", err)
					}
					go AddToDelete(sentMsg.Chat.ID, sentMsg.MessageID)
					log.Printf("Invoice %d paid!\n", invoice.ID)
					NewMessage(chatID, bot, fmt.Sprintf("Invoice %d paid!\n", invoice.ID), false)
					return
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func AddToDelete(chatID int64, messageID int) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	file, err := os.OpenFile("./messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error delete file: ",err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d:%d\n", chatID, messageID))
}
func ClearMessages(chatID int64, bot *tgbotapi.BotAPI) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Println("Error opening the delete file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var newData []string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		cid, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Printf("Error parsing chat ID: %s", err)
			continue
		}
		mid, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Error parsing message ID: %s", err)
			continue
		}
		if cid == chatID {
			bot.Request(tgbotapi.DeleteMessageConfig{
				ChatID:    cid,
				MessageID: mid,
			})
		} else {
			newData = append(newData, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from file: %s", err)
		return
	}

	file, err = os.Create("./messages.txt")
	if err != nil {
		log.Printf("Error opening the delete file for writing: %s", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newData {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Printf("Error writing to file: %s", err)
			return
		}
	}
	writer.Flush()
}