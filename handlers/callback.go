package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	switch callbackQuery.Data {
	//--------
	//---------
	//----------
	}
	return nil, nil, nil, nil
}
