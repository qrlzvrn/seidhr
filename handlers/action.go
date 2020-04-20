package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/qrlzvrn/seidhr/db"
	"github.com/qrlzvrn/seidhr/keyboards"
	"github.com/qrlzvrn/seidhr/med"
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
	isSubscribe, err := db.IsUserHasSub(conn, tguserID)
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

// SearchMed - меняет состояние пользователя на "SearchMed" и
// выдает конфиги сообщения, в котором пользователю предлагается
// ввести название лекарства
func SearchMed(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	tguserID := callbackQuery.From.ID
	err := db.ChangeUserState(conn, tguserID, "SearchMed")
	if err != nil {
		return nil, nil, nil, err
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboards.MedSearchKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Введите название лекарства:")

	return msg, newKeyboard, newText, nil
}

// SearchMedAct - принимает название лекарства и отправлят запрос на получение информации о нем,
// после возвращет конфиги сообщений
func SearchMedAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	medTitle := message.Text

	// Проверяем наличие данного лекарства в базе данных льготных лекарств
	isExist, err := db.IsMedExist(conn, medTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	trueName, err := db.FindTrueMedName(conn, medTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	if isExist == false {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, но кажется вы неправильно написали название, либо это лекарство не льготное. Попробуйте еще раз:")
		newKeyboard = nil
		newText = nil
	}

	// Отправляем запрос
	medResp, err := med.ReqMedInfo(trueName)
	if err != nil {
		return nil, nil, nil, err
	}

	// Находим id нашго лекарства
	medicomentID, err := db.FindMedID(conn, trueName)
	if err != nil {
		return nil, nil, nil, err
	}

	// Проверяем подписан ли пользователь на это лекарство, что бы решить
	// какую клавиатуру и текст необходимо отобразить
	isSubscribe, err := db.IsUserSubToThisMed(conn, tguserID, medicomentID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isSubscribe == true {
		// Проверяем полученный json на наличе информации об ошибке.
		// Так как перед отправкой запроса мы проверяем наличие лекарства в нашей бд,
		// где хранится список лекарств доступных по льготе, то вариант с неправильным написанием
		// или вводом чего-то вообще неподходящего или несуществующего
		// Значит, ошибка всегда будет означать то, что лекарства сейчас нет в доступе
		isErr := med.IsErrExistInJSON(medResp)
		if isErr == true {
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "К сожалению данного лекарства сейчас нет ни в одной аптеке, но так как вы подписаны, мы уведомим вас, как только оно появится в аптеках")
			msgConf.ReplyMarkup = keyboards.ViewMedKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil

			return msg, newKeyboard, newText, nil
		}

		// Парсим json и компануем текст сообщения
		msgText := med.ParseJSON(medResp)

		msgConf := tgbotapi.NewMessage(message.Chat.ID, msgText)
		msgConf.ReplyMarkup = keyboards.ViewMedKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		return msg, newKeyboard, newText, nil

	}

	// Опять Проверяем полученный json на наличе информации об ошибке
	isErr := med.IsErrExistInJSON(medResp)
	if isErr == true {

		msgConf := tgbotapi.NewMessage(message.Chat.ID, "К сожалению данного лекарства сейчас нет ни в одной аптеке, но если хотите, вы можете подписаться и мы уведомим вас, как только оно появится")
		msgConf.ReplyMarkup = keyboards.ViewMedWithSubKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
		return msg, newKeyboard, newText, nil
	}

	// Парсим json и компануем текст сообщения
	msgText := med.ParseJSON(medResp)

	msgConf := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msgConf.ReplyMarkup = keyboards.ViewMedWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// BackToHome - возвращает пользователя на домашнюю страницу
func BackToHome(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	tguserID := callbackQuery.From.ID

	isSubscribe, err := db.IsUserHasSub(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	db.ChangeUserState(conn, tguserID, "home")
	// Если у пользователя нет подписок, то выдаем ему клавиатуру
	// без кнопки просмотра подписок
	if isSubscribe == false {
		msg = nil

		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboards.HomeKeyboard)

		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Что бы вы хотели?")

		return msg, newKeyboard, newText, nil
	}

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboards.HomeWithSubKeyboard)

	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Что бы вы хотели?")

	return msg, newKeyboard, newText, nil
}

// Subscribe - оформляет подписку на лекарство для данного пользователя
func Subscribe(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	tguserID := callbackQuery.From.ID

	medTitle, err := db.CheckSelectedMed(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	medicamentID, err := db.FindMedID(conn, medTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.Subscribe(conn, tguserID, medicamentID)
	if err != nil {
		return nil, nil, nil, err
	}

	db.ChangeUserState(conn, tguserID, "home")

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboards.HomeWithSubKeyboard)

	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Поздравляю, подписка на лекарство успешно оформлена. Теперь вы первым узнаете о появлении данного лекарства в аптеках нашего города.\n\nХотите еще что-нибудь?")

	return msg, newKeyboard, newText, nil
}

// Unsubscribe - отменяет подписку на лекарство для данного пользователя
func Unsubscribe(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	tguserID := callbackQuery.From.ID

	medTitle, err := db.CheckSelectedMed(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	medicamentID, err := db.FindMedID(conn, medTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.Unsubscribe(conn, tguserID, medicamentID)
	if err != nil {
		return nil, nil, nil, err
	}

	db.ChangeUserState(conn, tguserID, "home")

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboards.HomeWithSubKeyboard)

	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Поздравляю, подписка отменена.\n\nХотите еще что-нибудь?")

	return msg, newKeyboard, newText, nil
}

// ListSubscriptions - отправляет пользователю сообщение с информацией о всех его подписках
func ListSubscriptions(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	tguserID := callbackQuery.From.ID

	subscriptions, err := db.ListUserSubscriptions(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	subKeyboard := keyboards.CreateKeyboarWithUserSubscriptions(subscriptions)

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, subKeyboard)

	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Ваши подписки:")

	return msg, newKeyboard, newText, nil
}
