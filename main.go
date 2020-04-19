package main

import (
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/seidhr/config"
	"github.com/qrlzvrn/seidhr/db"
	"github.com/qrlzvrn/seidhr/handlers"
	"github.com/qrlzvrn/seidhr/med"
)

// checkTime - следит за временем и в момент, когда время становится равным 11 часам пишет в канал
// Который будет считан функцией CyclicMedSearch, после чего она будет запущена.
func checkTime(c chan bool) {
	for {
		hour := time.Now().Hour()

		if hour == 11 {
			c <- true
			time.Sleep(23 * time.Hour)
		}
		time.Sleep(20 * time.Minute)
	}
}

func main() {

	conn, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// Проверяем заполнена ли наша база значениями полученными из файла drugs.txt
	// если нет, то заполняем.

	isExist, err := db.IsMedListExist(conn)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if isExist == false {

		meds, err := med.ReadFileWithMeds()
		if err != nil {
			log.Fatalf("%+v", err)
		}

		err = db.InitMedList(conn, meds)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}

	conf, err := config.InitConf()
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем бота
	if bot, err := tgbotapi.NewBotAPI(conf.TgBot.APIToken); err != nil {
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

		// Создаем канал необходимый для работы функций отвечающих за ежедневную проверку
		// наличия лекарств в аптеке.

		c := make(chan bool)

		go checkTime(c)
		go med.CyclicMedSearch(bot, c)

		//Начинаем слушать на 8443 порту
		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS(":8443", conf.SSL.Fullchain, conf.SSL.Privkey, nil)

		// Получаем обновления от телеграма,
		// в зависимости от типа полученного сообщения используем разные обработчики
		for update := range updates {
			if update.Message != nil {
				msg, newKeyboard, newText, err := handlers.MessageHandler(update.Message)
				if err != nil {
					log.Fatal(err)
				}
				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						log.Fatal(err)
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						log.Fatal(err)
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						log.Fatal(err)
					}
				}

			} else if update.CallbackQuery != nil {
				msg, newKeyboard, newText, err := handlers.CallbackHandler(update.CallbackQuery)
				if err != nil {
					log.Fatal(err)
				}
				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						log.Fatal(err)
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						log.Fatal(err)
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}
}
