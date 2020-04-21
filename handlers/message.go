package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/seidhr/db"
)

// MessageHandler - обрабатывает обычные сообщения от пользователей
// и в соответсвии с содержанием вызывает функции из action.go,
// после чего выдает три конфига для ответного сообщения
//
//		- msg - конфиг нового сообщение для пользователя
//		- newKeyboard - конфиг измененой клавиатуры
//		- newText - конфиг измененного текста сообщения
// ---------------------------------------------------------
// ВАЖНО:
// 			Данная функция занимается только анализом,
// 			все действия совершаются в action.go.
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	// Проверяем является полученное сообщение коммандой
	// или обычным текстовым сообщением и в зависимости от результата
	// обрабатываем определенным образом
	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			msg, newKeyboard, newText, err := Start(message)
			if err != nil {
				return nil, nil, nil, err
			}
			return msg, newKeyboard, newText, nil

		case "help":
			msg, newKeyboard, newText := Help(message)

			return msg, newKeyboard, newText, nil
		default:
			msg, newKeyboard, newText := Default(message)
			return msg, newKeyboard, newText, nil
		}
	} else {
		tguserID := message.From.ID

		conn, err := db.ConnectToDB()
		if err != nil {
			return nil, nil, nil, err
		}
		defer conn.Close()

		state, err := db.GetUserState(conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		switch state {
		// Производим запрос к Гос Услугам с помощью SearchMedAction() и
		// выводим
		case "SearchMed":
			msg, newKeyboard, newText, err := SearchMedAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "":
			////////////
			//////////
			////////
		}
	}
	return nil, nil, nil, nil
}
