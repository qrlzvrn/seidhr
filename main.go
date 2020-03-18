package main

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

//Функция заглушка
func runEverySecond() {
	fmt.Println("----")
}

func main() {
	// Заглушки для значений, которые позже будут получаться из конфига
	var botAPIToken string
	var fullchain string
	var privkey string

	// Используем cron для Go,
	// дабы по росписанию проверять наличие лекарств,
	// в случае, если лекарства появились, отправляем подписаным пользователям сообщения
	// -----------
	// Пока что используем runEverySecond как заглушку
	cronJob := cron.New()
	cronJob.Start()
	cronJob.AddFunc("* * * * * *", runEverySecond)

	// Инициализируем бота
	if bot, err := tgbotapi.NewBotAPI(botAPIToken); err != nil {
		log.Fatalf("%+v", err)
	} else {
		bot.Debug = true
		// Поулчаем инфу о состоянии нашего вебхука
		// Выводим в консоль последнюю возникшую ошибку
		info, err := bot.GetWebhookInfo()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}
		//Начинаем слушать на 8443 порту
		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS(":8443", fullchain, privkey, nil)

		// Получаем обновления от телеграма,
		// в зависимости от типа полученного сообщения используем разные обработчики
		for update := range updates {
			if update.Message != nil {
				//----------------
				//-----------------
				//------------------
			} else if update.CallbackQuery != nil {
				//----------------
				//-----------------
				//------------------
			}
		}
	}
	cronJob.Stop()
}
