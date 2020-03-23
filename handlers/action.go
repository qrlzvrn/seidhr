package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/seidhr/db"
	"github.com/qrlzvrn/seidhr/keyboards"
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
	tguserID := message.From.ID
	chatID := message.Chat.ID
	conn, err := db.ConnectToDB()
	if err != nil {
		return nil, nil, nil, err
	}
	defer conn.Close()

	isExist, err := db.CheckUser(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Проверяем первый ли раз пользователь обращается к боту
	if isExist == false {
		err := db.CreateUser(conn, tguserID, chatID)
		if err != nil {
			return nil, nil, nil, err
		}

		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Добро пожаловать.\n\nЧто бы проверить наличие необходимого вам льготного лекарства в аптеках Санкт-Петербурга, просто нажмите на кнопку и введите его навание.\n\nВ случае, если необходимого вам лекарства сейчас нигде нет, вы можете подписаться на него и мы сообщим вам, как только оно появится.\n\nДля получения информационной справки используте команду /help Приятного использования!")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
		return msg, newKeyboard, newText, nil
	}

	// Если пользователь уже взаимодействовал с ботом,
	// смотрим состояние его подписок
	isSubscribe, err := db.CheckSubscriptions(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Если у пользователя нет подписок, то выдаем ему клавиатуру
	// без кнопки просмотра подписок
	if isSubscribe == false {

		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Что бы вы хотели?")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
		return msg, newKeyboard, newText, nil
	}

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "Что бы вы хотели?")
	msgConf.ReplyMarkup = keyboards.HomeWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil

}

// Help ...
func Help(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	msg = tgbotapi.NewMessage(message.Chat.ID, "---Необольшая справка---\n\nС помощью данного бота вы можете проверить наличие льготных лекарств в аптеках Санкт-Петербурга, а так же подписаться на необходимые вам лекарства, и получать уведомления, как только они появятся в аптеках.\n\nЧто бы подписаться на какое-либо лекарство, вам необходимо нажать на кнопку <Проверить лекарство> и ввести его название, после чего в появившемся сообщении вы увидите всю информацию о нем, а так же кнопку <Подписаться>, если вы, конечно, уже не подписаны на него\n\nПосле того, как вы подписались на ваше первое лекарство, в главном меню появится кнопка <Подписки>, нажав на которую вы увидите все ваши подписки, узнать наличие, а так же отменить подписку.")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText
}

// Default ...
func Default(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, я так не умею :с")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText
}
