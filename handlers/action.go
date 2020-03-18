package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Данная часть пакета содержит в себе функции выступающие посредниками между
// пакетом db и функциями messageHandler и callbackHandler,
// а так же призвана повыстить читаемсть и расширяемость в будущем

// Обявим переменные для конфигов сообщений в самом начале,
// что бы не приходилось повторять их объявление в каждой функции
var msg, newKeyboard, newText tgbotapi.Chattable

// Реализация функционала команд /start, /help и тех, что будут добавлены в будущем

// Start ...
func Start(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	//------------
	//-------------
	//--------------
	return nil, nil, nil, nil
}

// Help ...
func Help(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	//------------
	//-------------
	//--------------
	return nil, nil, nil
}

// Default ...
func Default(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	//------------
	//-------------
	//--------------
	return nil, nil, nil
}
