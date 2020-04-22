package db

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/seidhr/config"
	"github.com/qrlzvrn/seidhr/errorz"
)

// ConnectToDB - подключается к базе данных и возвращаент конект
func ConnectToDB() (*sqlx.DB, error) {
	conf, err := config.InitConf()
	if err != nil {
		err := errorz.NewErrStack("Не удалось инициализировать конфиг для базы данных")
		return nil, err
	}

	port, err := strconv.Atoi(conf.DB.Port)
	if err != nil {
		err := errorz.NewErrStack("В конфигурационном файле порт базы данных указан некорректно")
		return nil, err
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.DB.Host, port, conf.DB.Username, conf.DB.Password, conf.DB.Name)

	if db, err := sqlx.Connect("postgres", dbInfo); err != nil {
		err := errorz.NewErrStack("Не удалось подключиться к базе данных")
		return nil, err
	} else {
		return db, nil
	}
}

// CreateUser - создает нового польователя
func CreateUser(db *sqlx.DB, tguserID int, chatID int64) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, chat_id, state, selected_med) VALUES ($1, $2, $3, $4)", tguserID, chatID, "born", 0)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("При попытке создания нового пользователя что-то пошло не так")
		return err
	}

	return nil
}

// CheckUser - проверяет наличие пользователя в базе
func CheckUser(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки существования пользователя закончился ошибкой")
		return false, err
	}
	return isExist, nil
}

// CheckUserState - проверяет состояние пользователя
func GetUserState(db *sqlx.DB, tguserID int) (string, error) {
	var state string
	err := db.QueryRow("SELECT state from tguser where id=$1", tguserID).Scan(&state)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе для получения состояния пользователя закончился ошибкой")
		return "", err
	}
	return state, nil
}

// ChangeUserState - изменяет состояние пользователя
func ChangeUserState(db *sqlx.DB, tguserID int, state string) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, tguserID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для изменения состояния пользователя закончился ошибкой")
		return err
	}
	return nil
}

// InitMedList - инициализирует список льготных лекарств в базе данных
func InitMedList(db *sqlx.DB, medLines []string) error {
	tx := db.MustBegin()

	for _, med := range medLines {
		tx.MustExec("INSERT INTO medicament (title, availability) VALUES ($1, $2)", med, false)
	}
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для заполнения таблицы medicament значениями закончился ошибкой")
		return err
	}
	return nil
}

// IsMedListExist - проверяет заполнена ли таблица medicament.
// Служит для того что бы не пытаться каждый раз заполнять бд значениями из файла drugs.txt
func IsMedListExist(db *sqlx.DB) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from medicament)").Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки существования списка лекарств закончился ошибкой")
		return false, err
	}
	return isExist, nil
}

// IsMedExist - проверяет существует ли такое лекарство в нашей базе
func IsMedExist(db *sqlx.DB, medName string) (bool, error) {

	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from medicament where title % $1)", medName).Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки существования лекарсвта закончился ошибкой")
		return false, err
	}
	return isExist, nil
}

// FindTrueMedName - выводит правильное название лекарства, если пользователь
// ввел название с опечатками.
//
// Данная функция используется в связке с IsMedExist.
//
//IsMedExist - проверяет существования лекарства, а данная функция
// выдает правильное название, для дальнешей работы с Гос. Услугами
func GetTrueMedName(db *sqlx.DB, medName string) (string, error) {

	var trueName string
	err := db.QueryRow("SELECT title FROM medicament WHERE title % $1", medName).Scan(&trueName)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения правильного названия лекарства закончился ошибкой")
		return "", err
	}
	return trueName, nil
}

// Subscribe - создает новую подписку пользователя на лекарство
func Subscribe(db *sqlx.DB, tguserID int, medicamentID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO subscription (tguser_id, medicament_id) VALUES ($1, $2)", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для оформления пописки закончился ошибкой")
		return err
	}
	return nil
}

// Unsubscribe - отменяет у пользьователя подписку на лекарство
func Unsubscribe(db *sqlx.DB, tguserID int, medicamentID int) error {
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM subscription WHERE tguser_id = $1 AND medicament_id = $2", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для отмены подписки закончился ошибкой")
		return err
	}
	return nil
}

// GetUserSubscriptions - находит все подписки пользователя и возвращает [][]string, где
// [[id title] [id title] [id title]]
func GetUserSubscriptions(db *sqlx.DB, tguserID int) ([][]string, error) {
	rows, err := db.Query("SELECT id, title from medicament INNER JOIN subscription on medicament.id=subscription.medicament_id WHERE subscription.tguser_id = $1", tguserID)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения всех подписок пользователя закончился ошибкой")
		return nil, err
	}

	subscriptions := [][]string{}

	for rows.Next() {
		var id int
		var title string

		rows.Scan(&id, &title)
		subscriptions = append(subscriptions, []string{strconv.Itoa(id), title})
		defer rows.Close()
	}
	return subscriptions, nil
}

// IsUserSubToThisMed - проверяет подписан ли пользователь на данное лекарство
func IsUserSubToThisMed(db *sqlx.DB, tguserID int, medicamentID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1 AND medicament_id=$2)", tguserID, medicamentID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки того, подписан ли пользователь на данное лекарство закончился ошибкой")
		return false, err
	}
	return isExist, nil
}

// IsUserHasSub - проверяет наличие у пользователя подписок на лекарства и,
// если у него есть хоть одна подписка, возвращает true.
func IsUserHasSub(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки наличия у пользователя хотя бы одной подписки закончился ошибкой")
		return false, err
	}
	return isExist, nil
}

// GetAllMedicamentsWithSub - Находит все лекарства, на которые подписаны пользователи
// и возварщает слайс с их id
func GetAllMedicamentsWithSub(db *sqlx.DB) ([]int, error) {
	rows, err := db.Query("SELECT DISTINCT medicament_id FROM subscription")
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для нахождения всех лекарств с оформленными подписками закончился ошибкой")
		return nil, err
	}

	subMeds := []int{}

	for rows.Next() {
		var id int

		rows.Scan(&id)
		subMeds = append(subMeds, id)
		defer rows.Close()
	}
	return subMeds, nil
}

// GetSubscribers - находит пользователей подписанных на определенное лекарство,
// id которого принимается на вход, и возвращает слайс с id этих пользователей
func GetSubscribers(db *sqlx.DB, medicamentID int) ([]int, error) {
	rows, err := db.Query("SELECT tguser_id FROM subscription WHERE medicament_id = $1", medicamentID)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения всех пользователей подписанных на данное лекарство закончился ошибкой")
		return nil, err
	}

	users := []int{}

	for rows.Next() {
		var id int

		rows.Scan(&id)
		users = append(users, id)
		defer rows.Close()
	}
	return users, nil
}

// FindChatID - находит пользователя и возвращает его chatID
func GetChatID(db *sqlx.DB, tguserID int) (int, error) {
	var chatID int
	err := db.QueryRow("SELECT chat_id FROM tguser WHERE id = $1", tguserID).Scan(&chatID)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения chatID пользователя закончился ошибкой")
		return 0, err
	}
	return chatID, nil
}

// CheckAvailability - проверяет наличие лекарства записаное в базе
func GetAvailability(db *sqlx.DB, medicamentID int) (bool, error) {
	var availible bool
	err := db.QueryRow("SELECT availability FROM medicament WHERE id = $1", medicamentID).Scan(&availible)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки наличия лекарства закончился ошибкой")
		return false, err
	}
	return availible, nil
}

// ChangeAvailability - изменяет наличие лекарства в базе
func ChangeAvailability(db *sqlx.DB, medicamentID int, value bool) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE medicament SET availability = $1 WHERE id = $2", value, medicamentID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для изменения доступности лекарства закончился ошибкой")
		return err
	}
	return nil
}

// FindMedID - находит id необхомодимого лекартсва
func GetMedID(db *sqlx.DB, medTitle string) (int, error) {
	var medicamentID int
	err := db.QueryRow("SELECT id FROM medicament WHERE title = $1", medTitle).Scan(&medicamentID)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения id лекарства закончился ошибкой")
		return 0, err
	}
	return medicamentID, nil
}

// FindMedTitle - находит название лекарства по его id
func GetMedTitle(db *sqlx.DB, medicamentID int) (string, error) {
	var medTitle string
	err := db.QueryRow("SELECT title FROM medicament WHERE id = $1", medicamentID).Scan(&medTitle)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения названия лекарства закончился ошибкой")
		return "", err
	}
	return medTitle, nil
}

// AreTheAnySubscriptions - проверяет существование хотя бы одной подписки
// Служит для того, что бы избежать ошибок в функции CyclicMedSearch
func AreTheAnySubscriptions(db *sqlx.DB) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM subscription )").Scan(&isExist)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для проверки существования хотя бы одной подписки закончился ошибкой")
		return false, err
	}

	return isExist, nil
}

// CheckSelectedMed - получает id лекарства выбранного пользователем в данный момент
func GetSelectedMed(db *sqlx.DB, tguserID int) (int, error) {
	var medicamentID int
	err := db.QueryRow("SELECT selected_med FROM tguser WHERE id = $1", tguserID).Scan(&medicamentID)
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для получения выбранного пользователем лекарства закончился ошибкой")
		return 0, err
	}
	return medicamentID, nil
}

// ChangeSelectedMed - меняет выбранное пользователем лекарство
func ChangeSelectedMed(db *sqlx.DB, medicamentID, tguserID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET selected_med = $1 WHERE id = $2", medicamentID, tguserID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErrStack("Запрос к базе данных для изменения выбранного пользователем лекарства закончился ошибкой")
		return err
	}
	return nil
}
