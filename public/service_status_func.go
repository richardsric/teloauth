package public

import (
	"log"
	"strconv"

	api "gopkg.in/telegram-bot-api.v4"
)

//SendServiceStatusIM this use to send HTML parsed  message to a telegram user.
func SendServiceStatusIM(msg string) {
	//BotKey for Error Reporting
	var adminTelegram = "430073910"
	var botKey = "482344425:AAHVf37OosNE7iBKLEYRG-wN8raNs_lk2xY"

	bot, err := api.NewBotAPI(botKey)
	if err != nil {
		log.Println("SendServiceStatusIM", err)
	}
	bot.Debug = false
	cID, _ := strconv.ParseInt(adminTelegram, 10, 64)
	ms := api.NewMessage(cID, msg)
	ms.ParseMode = "HTML"

	bot.Send(ms)
}
