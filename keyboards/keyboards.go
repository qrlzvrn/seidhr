package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// HomeKeyboard - клавиатура домашнего экрана, пользователй,
// у которых еще нет ни одной подписки
var HomeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Проверить лекарство", "searchMed"),
	),
)

// HomeWithSubKeyboard - клавиатура домашнего экрана,
// с возможностью просмотра подписок
var HomeWithSubKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Проверить лекарство", "searchMed"),
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
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отписаться", "unsubscribe"),
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

// CreateKeyboarWithUserSubscriptions - генерирует клавиатуру со списком подписок.
// На кнопках отображаются названия лекарств, нк которые оформлены подписки,
// а в качестве data берутся их id
func CreateKeyboarWithUserSubscriptions(subscriptions [][]string) tgbotapi.InlineKeyboardMarkup {
	AllSubscriptionsKeyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, sub := range subscriptions {
		id := sub[0]
		name := sub[1]

		var row []tgbotapi.InlineKeyboardButton

		btn := tgbotapi.NewInlineKeyboardButtonData(name, id)
		row = append(row, btn)
		AllSubscriptionsKeyboard.InlineKeyboard = append(AllSubscriptionsKeyboard.InlineKeyboard, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData("Назад", "backToHome")
	row = append(row, btn)
	AllSubscriptionsKeyboard.InlineKeyboard = append(AllSubscriptionsKeyboard.InlineKeyboard, row)
	return AllSubscriptionsKeyboard
}
