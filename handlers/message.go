package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
			//-------
			//--------
			//---------
			return nil, nil, nil, nil

		case "help":
			//-------
			//--------
			//---------
			return nil, nil, nil, nil
		default:
			//-------
			//--------
			//---------
			return nil, nil, nil, nil
		}
	} else {
		//-------
		//--------
		//---------
	}
	return nil, nil, nil, nil
}
