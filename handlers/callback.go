package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/seidhr/db"
)

// CallbackHandler - обрабатывает сообщения от нажатий на inlineKeyboard
// и в ависимости от поля Data этой кнопки, вызывает функцию из action.go,
// после чего выдает три конфига для ответного сообщения
//
//		- msg - конфиг нового сообщение для пользователя
//		- newKeyboard - конфиг измененой клавиатуры
//		- newText - конфиг измененного текста сообщения
// ---------------------------------------------------------
// ВАЖНО:
// 			Данная функция занимается только анализом,
// 			все действия совершаются в action.go.
func CallbackHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	conn, err := db.ConnectToDB()
	if err != nil {
		return nil, nil, nil, err
	}

	switch callbackQuery.Data {
	//Перехватываем нажатие на кнопку <Проверить лекарство> и
	// с помощью SearchMed() меняем state пользователя на "SearchMed"
	// и предлагаем ввести название лекарства
	case "searchMed":
		msg, newKeyboard, newText, err := SearchMed(callbackQuery, conn)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil

	case "backToHome":
		msg, newKeyboard, newText, err := BackToHome(callbackQuery, conn)
		if err != nil {
			return nil, nil, nil, err
		}
		return msg, newKeyboard, newText, nil

	case "subscribe":
		msg, newKeyboard, newText, err := Subscribe(callbackQuery, conn)
		if err != nil {
			return nil, nil, nil, err
		}
		return msg, newKeyboard, newText, nil

	case "usnubscribe":
		msg, newKeyboard, newText, err := Unsubscribe(callbackQuery, conn)
		if err != nil {
			return nil, nil, nil, err
		}
		return msg, newKeyboard, newText, nil
	}
	return msg, newKeyboard, newText, nil
}
