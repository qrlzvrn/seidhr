package db

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/seidhr/config"
)

// ConnectToDB ...
func ConnectToDB() (*sqlx.DB, error) {
	dbConf, err := config.NewDBConf()
	if err != nil {
		return nil, err
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Name)

	if db, err := sqlx.Connect("postgres", dbInfo); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

// CreateUser ...
func CreateUser(db *sqlx.DB, tguserID int, chatID int64) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, chat_id) VALUES ($1, $2)", tguserID, chatID)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CheckUser ...
func CheckUser(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// InitMedList ...
func InitMedList(db sqlx.DB, medLines []string) error {
	tx := db.MustBegin()

	for _, med := range medLines {
		tx.MustExec("INSERT INTO medicament (title) VALUES ($1)", med)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// CheckMed ...
func CheckMed(db sqlx.DB, medName string) (bool, error) {

	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from medicament where title=$1)", medName).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// Subscribe - создает новую подписку пользователя на лекарство
func Subscribe(db sqlx.DB, tguserID int, medicamentID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO subscription (tguser_id, medicament_id) VALUES ($1, $2)", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
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
		return err
	}
	return nil
}

// ListSubscriptions - находит все подписки пользователя и возвращает [][]string, где
// [[id title] [id title] [id title]]
func ListSubscriptions(db *sqlx.DB, tguserID int) ([][]string, error) {
	rows, err := db.Query("SELECT id, title from medicament INNER JOIN subscription on medicament.id=subscription.medicament_id WHERE subscription.tguser_id = $1", tguserID)
	if err != nil {
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

// CheckSubscriptions - проверяет наличие у пользователя подписок на лекарства и,
// если у него есть хоть одна подписка, возвращает true.
func CheckSubscriptions(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// FindSubMed - Находит все лекарства, на которые подписаны пользователи
// и возварщает слайс с их id
func FindSubMed(db *sqlx.DB) ([]int, error) {
	rows, err := db.Query("SELECT DISTINCT medicament_id FROM subscription")
	if err != nil {
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

// FindUsersWhoSub - находит пользователей подписанных на определенное лекарство,
// id которого принимается на вход, и возвращает слайс с id этих пользователей
func FindUsersWhoSub(db *sqlx.DB, medicamentID int) ([]int, error) {
	rows, err := db.Query("SELECT tguser_id FROM subscription WHERE medicament_id = $1", medicamentID)
	if err != nil {
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
func FindChatID(db *sqlx.DB, tguserID int) (int, error) {
	var chatID int
	err := db.QueryRow("SELECT chat_id FROM tguser WHERE id = $1", tguserID).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func FindMed(db *sqlx.DB, medicamentID int) (string, error) {
	var medTitle string
	err := db.QueryRow("SELECT title FROM medicament WHERE id = $1", medicamentID).Scan(&medTitle)
	if err != nil {
		return "", err
	}
	return medTitle, nil
}

func CheckAvailability(db *sqlx.DB, medicamentID int) (bool, error) {
	var availible bool
	err := db.QueryRow("SELECT availability FROM medicament WHERE id = $1", medicamentID).Scan(&availible)
	if err != nil {
		return false, err
	}
	return availible, nil
}

func ChangeAvailability(db *sqlx.DB, medicamentID int, value bool) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE medicament SET availability = $1 WHERE id = $2", value, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
