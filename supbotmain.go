package main

import(

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
				SupStartMenu()
			} 
		}
		if update.CallbackQuery != nil {
			upCQ := update.CallbackQuery;
			switch upCQ.Data {
			case "":
			}
		}
	}
}
