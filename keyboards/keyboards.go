package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// HomeKeyboard - клавиатура домашнего экрана, пользователй,
// у которых еще нет ни одной подписки
var HomeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Проверить лекарство", "checkMed"),
	),
)

// HomeWithSubKeyboard - клавиатура домашнего экрана,
// с возможностью просмотра подписок
var HomeWithSubKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Проверить лекарство", "checkMed"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подписки", "lsSub"),
	),
)

// MedSearchKeyboard - клавиатура при запросе ввода названия лекарства
var MedSearchKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToHome"),
	),
)

// ViewMedKeyboard - клавиатура просмотра лекарства
var ViewMedKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToHome"),
	),
)

// ViewMedWithSubKeyboard - клавиатура просмотра лекарства
// с возможностью подписаться на него
var ViewMedWithSubKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Подписаться", "subscribe"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToHome"),
	),
)

// ListSubsKeyboard - клавиатура просмотра подписок
var ListSubsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToHome"),
	),
)
